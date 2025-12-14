// Package main provides a CLI tool to remove comments from source code.
package main

import (
	"os"

	"github.com/connerohnesorge/rm-comments/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
