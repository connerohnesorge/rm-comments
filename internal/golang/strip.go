package golang

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// StripOptions configures comment removal behavior.
type StripOptions struct {
	// RemoveDirectives removes compiler directives (//go:, // +build)
	// when true. Otherwise they are preserved.
	RemoveDirectives bool
}

// RemoveComments removes comments from Go source code using AST.
// When opts.RemoveDirectives is true, compiler directives are also removed.
func RemoveComments(src []byte, removeDirectives bool) ([]byte, error) {
	opts := StripOptions{RemoveDirectives: removeDirectives}

	return StripComments(src, opts)
}

// StripComments removes comments from Go source code using AST.
// Compiler directives are preserved unless opts.RemoveDirectives is true.
func StripComments(src []byte, opts StripOptions) ([]byte, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if len(file.Comments) == 0 {
		return src, nil
	}

	preserveRanges := buildPreserveRanges(file.Comments, opts)
	result := removeCommentsFromSource(src, file.Comments, preserveRanges, fset)

	formatted, err := format.Source(result)
	if err != nil {
		// If formatting fails, return the unformatted result.
		// This shouldn't happen for valid Go code.
		return result, nil
	}

	return formatted, nil
}

// buildPreserveRanges builds a map of comment positions to preserve.
func buildPreserveRanges(
	comments []*ast.CommentGroup,
	opts StripOptions,
) map[token.Pos]token.Pos {
	preserve := make(map[token.Pos]token.Pos)

	if opts.RemoveDirectives {
		return preserve
	}

	for _, cg := range comments {
		for _, c := range cg.List {
			if IsDirective(c.Text) {
				preserve[c.Pos()] = c.End()
			}
		}
	}

	return preserve
}

// commentRange represents a byte range for a comment in source code.
type commentRange struct {
	start, end int
}

// removeCommentsFromSource removes comments from source code.
func removeCommentsFromSource(
	src []byte,
	commentGroups []*ast.CommentGroup,
	preserve map[token.Pos]token.Pos,
	fset *token.FileSet,
) []byte {
	ranges := collectCommentRanges(commentGroups, preserve, fset)

	if len(ranges) == 0 {
		return src
	}

	return buildResultWithoutComments(src, ranges)
}

// collectCommentRanges collects the byte ranges of comments to remove.
func collectCommentRanges(
	commentGroups []*ast.CommentGroup,
	preserve map[token.Pos]token.Pos,
	fset *token.FileSet,
) []commentRange {
	var ranges []commentRange

	for _, cg := range commentGroups {
		for _, c := range cg.List {
			if _, preserved := preserve[c.Pos()]; preserved {
				continue
			}

			start := fset.Position(c.Pos()).Offset
			end := fset.Position(c.End()).Offset
			ranges = append(ranges, commentRange{start, end})
		}
	}

	return ranges
}

// buildResultWithoutComments builds source without the specified ranges.
func buildResultWithoutComments(src []byte, ranges []commentRange) []byte {
	var result []byte
	lastEnd := 0

	for _, r := range ranges {
		result = append(result, src[lastEnd:r.start]...)
		lastEnd = r.end
	}

	result = append(result, src[lastEnd:]...)

	return result
}

// IsDirective checks if a comment is a Go compiler directive.
func IsDirective(text string) bool {
	if strings.HasPrefix(text, "//go:") {
		return true
	}

	if strings.HasPrefix(text, "// +build") {
		return true
	}

	if strings.HasPrefix(text, "//+build") {
		return true
	}

	return false
}
