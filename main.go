package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
	ExitMisuse  = 2
)

// Default directories to exclude during recursive traversal
var defaultExcludeDirs = []string{
	".git",
	"vendor",
	"node_modules",
	".direnv",
	"__pycache__",
	"target",
}

// GlobalFlags contains flags shared across all subcommands
type GlobalFlags struct {
	Recursive        bool     `short:"r" help:"Process directories recursively."`
	DryRun           bool     `short:"n" help:"Preview changes without modifying files."`
	Verbose          bool     `short:"v" help:"Print path of each modified file."`
	Include          []string `help:"Include only files matching glob pattern (can be repeated)."`
	Exclude          []string `help:"Exclude files matching glob pattern (can be repeated)."`
	RemoveDirectives bool     `help:"Remove compiler directives that would otherwise be preserved."`
}

// GoCmd handles Go comment removal
type GoCmd struct {
	Paths []string `arg:"" name:"path" help:"Paths to process. Use '-' for stdin/stdout." type:"path"`
}

// CLI is the root command structure
var CLI struct {
	GlobalFlags

	Go GoCmd `cmd:"" help:"Remove comments from Go files."`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("rm-comments"),
		kong.Description("Remove comments from source code files."),
		kong.UsageOnError(),
		kong.Exit(func(code int) {
			// Kong uses various exit codes; normalize to ExitMisuse for CLI errors
			if code != 0 {
				os.Exit(ExitMisuse)
			}
			os.Exit(code)
		}),
	)

	err := run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(ExitError)
	}
}

func run(ctx *kong.Context) error {
	switch ctx.Command() {
	case "go <path>":
		return runGo(CLI.Go.Paths, CLI.GlobalFlags)
	default:
		return fmt.Errorf("unknown command: %s", ctx.Command())
	}
}

// runGo processes Go files to remove comments
func runGo(paths []string, flags GlobalFlags) error {
	for _, path := range paths {
		if path == "-" {
			return processStdin(flags)
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot access %s: %w", path, err)
		}

		if info.IsDir() {
			if !flags.Recursive {
				return fmt.Errorf("cannot process directory %s without -r flag", path)
			}
			if err := processDirectory(path, flags); err != nil {
				return err
			}
		} else {
			if err := processFile(path, flags); err != nil {
				return err
			}
		}
	}
	return nil
}

// processStdin reads from stdin, removes comments, and writes to stdout
func processStdin(flags GlobalFlags) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	output, err := removeComments(input, flags.RemoveDirectives)
	if err != nil {
		return fmt.Errorf("processing stdin: %w", err)
	}

	_, err = os.Stdout.Write(output)
	return err
}

// processDirectory recursively processes all .go files in a directory
func processDirectory(dir string, flags GlobalFlags) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if d.IsDir() {
			name := d.Name()
			for _, excluded := range defaultExcludeDirs {
				if name == excluded {
					return filepath.SkipDir
				}
			}
			// Check custom exclude patterns on directory
			if shouldExclude(name, flags.Exclude) {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Check include/exclude patterns
		base := filepath.Base(path)
		if len(flags.Include) > 0 && !shouldInclude(base, flags.Include) {
			return nil
		}
		if shouldExclude(base, flags.Exclude) {
			return nil
		}

		return processFile(path, flags)
	})
}

// shouldInclude checks if filename matches any include pattern
func shouldInclude(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// shouldExclude checks if filename matches any exclude pattern
func shouldExclude(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// processFile removes comments from a single Go file
func processFile(path string, flags GlobalFlags) error {
	// Verify .go extension
	if !strings.HasSuffix(path, ".go") {
		return fmt.Errorf("not a Go file: %s", path)
	}

	// Read original file
	original, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	// Remove comments
	modified, err := removeComments(original, flags.RemoveDirectives)
	if err != nil {
		return fmt.Errorf("processing %s: %w", path, err)
	}

	// Skip if no changes
	if bytes.Equal(original, modified) {
		return nil
	}

	// Dry-run mode: show diff
	if flags.DryRun {
		diff := generateDiff(path, original, modified)
		fmt.Print(diff)
		return nil
	}

	// Verbose mode: print path
	if flags.Verbose {
		fmt.Println(path)
	}

	// Write modified content, preserving permissions
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	return os.WriteFile(path, modified, info.Mode().Perm())
}

// removeComments removes comments from Go source code using AST
func removeComments(src []byte, removeDirectives bool) ([]byte, error) {
	fset := token.NewFileSet()

	// Parse with comments to identify them
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// If no comments, return original
	if len(file.Comments) == 0 {
		return src, nil
	}

	// Build a set of comment positions to preserve (directives)
	preserveRanges := make(map[token.Pos]token.Pos)
	if !removeDirectives {
		for _, cg := range file.Comments {
			for _, c := range cg.List {
				if isDirective(c.Text) {
					preserveRanges[c.Pos()] = c.End()
				}
			}
		}
	}

	// Rebuild the source without non-preserved comments
	result := removeCommentsFromSource(src, file.Comments, preserveRanges, fset)

	// Format the output using gofmt
	formatted, err := format.Source(result)
	if err != nil {
		// If formatting fails, return the unformatted result
		// This shouldn't happen for valid Go code
		return result, nil
	}

	return formatted, nil
}

// removeCommentsFromSource removes comments from source code
func removeCommentsFromSource(src []byte, commentGroups []*ast.CommentGroup, preserve map[token.Pos]token.Pos, fset *token.FileSet) []byte {
	// Collect all comment ranges to remove
	type commentRange struct {
		start, end int
	}
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

	if len(ranges) == 0 {
		return src
	}

	// Build result by copying non-comment parts
	var result []byte
	lastEnd := 0

	for _, r := range ranges {
		// Copy content before this comment
		result = append(result, src[lastEnd:r.start]...)
		lastEnd = r.end
	}

	// Copy remaining content after last comment
	result = append(result, src[lastEnd:]...)

	return result
}

// isDirective checks if a comment is a Go compiler directive
func isDirective(text string) bool {
	// Check for //go: directives
	if strings.HasPrefix(text, "//go:") {
		return true
	}

	// Check for legacy // +build constraints
	if strings.HasPrefix(text, "// +build") || strings.HasPrefix(text, "//+build") {
		return true
	}

	return false
}

// generateDiff creates a unified diff-style output
func generateDiff(path string, original, modified []byte) string {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("--- %s\n", path))
	buf.WriteString(fmt.Sprintf("+++ %s\n", path))

	origLines := strings.Split(string(original), "\n")
	modLines := strings.Split(string(modified), "\n")

	// Simple line-by-line diff
	maxLines := len(origLines)
	if len(modLines) > maxLines {
		maxLines = len(modLines)
	}

	inHunk := false
	hunkStart := 0

	for i := 0; i < maxLines; i++ {
		origLine := ""
		modLine := ""
		if i < len(origLines) {
			origLine = origLines[i]
		}
		if i < len(modLines) {
			modLine = modLines[i]
		}

		if origLine != modLine {
			if !inHunk {
				hunkStart = i + 1
				buf.WriteString(fmt.Sprintf("@@ -%d +%d @@\n", hunkStart, hunkStart))
				inHunk = true
			}
			if i < len(origLines) && origLine != "" {
				buf.WriteString(fmt.Sprintf("-%s\n", origLine))
			}
			if i < len(modLines) && modLine != "" {
				buf.WriteString(fmt.Sprintf("+%s\n", modLine))
			}
		} else if inHunk {
			// Context line
			buf.WriteString(fmt.Sprintf(" %s\n", origLine))
		}
	}

	return buf.String()
}
