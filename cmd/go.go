// Package cmd provides CLI command implementations for rm-comments.
// This package contains the root CLI structure and subcommands for
// processing source code files to remove comments.
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/connerohnesorge/rm-comments/internal/golang"
)

// GoCmd handles Go comment removal.
// It processes one or more paths, which can be files, directories,
// or "-" for stdin/stdout mode.
type GoCmd struct {
	Paths []string `arg:"" name:"path" type:"path" help:"Paths; '-' for stdin."`
}

// Run processes Go files to remove comments.
// For each path, it determines whether to process as stdin, directory,
// or single file. Directories require the -r flag.
func (c *GoCmd) Run(flags GlobalFlags) error {
	for _, path := range c.Paths {
		if path == "-" {
			return processStdin(flags)
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot access %s: %w", path, err)
		}

		if info.IsDir() {
			if !flags.Recursive {
				return fmt.Errorf(
					"cannot process directory %s without -r flag", path,
				)
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

// processStdin reads Go source from stdin, removes comments,
// and writes the result to stdout.
func processStdin(flags GlobalFlags) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	output, err := golang.RemoveComments(input, flags.RemoveDirectives)
	if err != nil {
		return fmt.Errorf("processing stdin: %w", err)
	}

	_, err = os.Stdout.Write(output)

	return err
}

// processDirectory recursively processes all .go files in a directory.
// It walks the directory tree and processes each Go file found.
func processDirectory(dir string, flags GlobalFlags) error {
	return filepath.WalkDir(dir, processEntry(flags))
}

// processEntry returns a WalkDirFunc that processes directory entries.
// It handles directory filtering and delegates file processing.
func processEntry(flags GlobalFlags) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return handleDirectory(d.Name(), flags)
		}

		return handleFile(path, flags)
	}
}

// handleDirectory checks if a directory should be skipped.
// It skips default excluded directories (like .git, vendor) and
// directories matching custom exclude patterns.
func handleDirectory(name string, flags GlobalFlags) error {
	if slices.Contains(DefaultExcludeDirs, name) {
		return filepath.SkipDir
	}

	if shouldExclude(name, flags.Exclude) {
		return filepath.SkipDir
	}

	return nil
}

// handleFile processes a file if it matches the filter criteria.
// Non-Go files are skipped. Include/exclude patterns are applied.
func handleFile(path string, flags GlobalFlags) error {
	if !strings.HasSuffix(path, ".go") {
		return nil
	}

	base := filepath.Base(path)
	if len(flags.Include) > 0 && !shouldInclude(base, flags.Include) {
		return nil
	}

	if shouldExclude(base, flags.Exclude) {
		return nil
	}

	return processFile(path, flags)
}

// shouldInclude checks if filename matches any include pattern.
// Returns true if the name matches at least one pattern.
func shouldInclude(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// shouldExclude checks if filename matches any exclude pattern.
// Returns true if the name matches at least one pattern.
func shouldExclude(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err == nil && matched {
			return true
		}
	}

	return false
}

// processFile removes comments from a single Go file.
// In dry-run mode, it prints a diff. Otherwise, it writes the
// modified content back to the file, preserving permissions.
func processFile(path string, flags GlobalFlags) error {
	if !strings.HasSuffix(path, ".go") {
		return fmt.Errorf("not a Go file: %s", path)
	}

	original, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}

	modified, err := golang.RemoveComments(original, flags.RemoveDirectives)
	if err != nil {
		return fmt.Errorf("processing %s: %w", path, err)
	}

	if bytes.Equal(original, modified) {
		return nil
	}

	if flags.DryRun {
		diff := golang.GenerateDiff(path, original, modified)
		fmt.Print(diff)

		return nil
	}

	if flags.Verbose {
		fmt.Println(path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	return os.WriteFile(path, modified, info.Mode().Perm())
}
