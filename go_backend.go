package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
)

func processGoFile(filename string, content []byte) (string, error) {
	fset := token.NewFileSet()
	
	// Parse the Go source file without comments
	// By not using parser.ParseComments mode, comments are not included in the AST
	file, err := parser.ParseFile(fset, filename, content, 0)
	if err != nil {
		return "", fmt.Errorf("failed to parse Go file: %w", err)
	}

	// Print the AST back to source code
	var buf bytes.Buffer
	cfg := printer.Config{
		Mode:     printer.TabIndent,
		Tabwidth: 8,
	}
	
	if err := cfg.Fprint(&buf, fset, file); err != nil {
		return "", fmt.Errorf("failed to print Go AST: %w", err)
	}

	return buf.String(), nil
}
