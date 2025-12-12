package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Supported extensions: .go, .ts, .tsx, .js, .jsx, .py, .rs\n")
		os.Exit(1)
	}

	filename := os.Args[1]
	ext := strings.ToLower(filepath.Ext(filename))

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var result string
	var processErr error

	switch ext {
	case ".go":
		result, processErr = processGoFile(filename, content)
	case ".ts", ".tsx", ".js", ".jsx":
		result, processErr = processTypeScriptFile(filename, content, ext)
	case ".py":
		result, processErr = processPythonFile(filename, content)
	case ".rs":
		result, processErr = processRustFile(filename, content)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported file extension: %s\n", ext)
		fmt.Fprintf(os.Stderr, "Supported extensions: .go, .ts, .tsx, .js, .jsx, .py, .rs\n")
		os.Exit(1)
	}

	if processErr != nil {
		fmt.Fprintf(os.Stderr, "Error processing file: %v\n", processErr)
		os.Exit(1)
	}

	fmt.Print(result)
}
