// Package golang provides Go source code manipulation utilities.
// This package contains functions for removing comments from Go source
// and generating unified diffs of the changes.
package golang

import (
	"fmt"
	"strings"
)

// GenerateDiff creates a unified diff-style output comparing original
// and modified content. The output format follows standard unified diff
// conventions with --- and +++ headers, and @@ hunk markers.
func GenerateDiff(path string, original, modified []byte) string {
	var buf strings.Builder

	fmt.Fprintf(&buf, "--- %s\n", path)
	fmt.Fprintf(&buf, "+++ %s\n", path)

	origLines := strings.Split(string(original), "\n")
	modLines := strings.Split(string(modified), "\n")
	maxLines := max(len(modLines), len(origLines))

	writeDiffLines(&buf, origLines, modLines, maxLines)

	return buf.String()
}

// writeDiffLines writes the diff content for each line.
// It iterates through both line slices and outputs hunks where they differ.
func writeDiffLines(
	buf *strings.Builder,
	origLines, modLines []string,
	maxLines int,
) {
	inHunk := false

	for i := range maxLines {
		origLine := getLine(origLines, i)
		modLine := getLine(modLines, i)

		if origLine != modLine {
			if !inHunk {
				fmt.Fprintf(buf, "@@ -%d +%d @@\n", i+1, i+1)
				inHunk = true
			}

			writeLinePair(buf, origLine, modLine)
		} else if inHunk {
			fmt.Fprintf(buf, " %s\n", origLine)
		}
	}
}

// getLine returns the line at index i, or empty string if out of bounds.
// This safely handles cases where the two files have different line counts.
func getLine(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}

	return ""
}

// writeLinePair writes a removed line (if non-empty) and an added line
// (if non-empty) to the diff output.
func writeLinePair(buf *strings.Builder, origLine, modLine string) {
	if origLine != "" {
		fmt.Fprintf(buf, "-%s\n", origLine)
	}

	if modLine != "" {
		fmt.Fprintf(buf, "+%s\n", modLine)
	}
}
