package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

// Exit codes for the application.
const (
	ExitSuccess = 0
	ExitError   = 1
	ExitMisuse  = 2
)

// DefaultExcludeDirs contains directories to exclude during recursive
// traversal.
var DefaultExcludeDirs = []string{
	".git",
	"vendor",
	"node_modules",
	".direnv",
	"__pycache__",
	"target",
}

// GlobalFlags contains flags shared across all subcommands.
type GlobalFlags struct {
	Recursive        bool     `short:"r" help:"Recurse directories."`
	DryRun           bool     `short:"n" help:"Preview only, no write."`
	Verbose          bool     `short:"v" help:"Print modified paths."`
	Include          []string `help:"Include glob (repeatable)."`
	Exclude          []string `help:"Exclude glob (repeatable)."`
	RemoveDirectives bool     `help:"Remove directives too."`
}

// CLI is the root command structure.
type CLI struct {
	GlobalFlags

	Go GoCmd `cmd:"" help:"Remove comments from Go files."`
}

// Run parses CLI arguments and executes the appropriate command.
// Returns an exit code suitable for os.Exit.
func Run() int {
	var cli CLI

	ctx := kong.Parse(&cli,
		kong.Name("rm-comments"),
		kong.Description("Remove comments from source code files."),
		kong.UsageOnError(),
		kong.Exit(func(_ int) {
			// Exit handled by returning exit code from Run
		}),
	)

	err := execute(ctx, &cli)
	if err != nil {
		_, _ = fmt.Fprintf(ctx.Stderr, "error: %v\n", err)

		return ExitError
	}

	return ExitSuccess
}

func execute(ctx *kong.Context, cli *CLI) error {
	switch ctx.Command() {
	case "go <path>":
		return cli.Go.Run(cli.GlobalFlags)
	default:
		return fmt.Errorf("unknown command: %s", ctx.Command())
	}
}
