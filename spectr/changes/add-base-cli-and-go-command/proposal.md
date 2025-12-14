# Change: Add Base CLI Structure and Go Comment Removal Command

## Why
The project currently has only a placeholder CLI skeleton with example commands (`rm`, `ls`) that don't match the project's purpose. We need to establish the foundational CLI structure with Kong and implement the first language-specific subcommand for Go comment removal using Go's native `go/ast` and `go/parser` packages.

## What Changes
- **BREAKING**: Replace placeholder CLI structure with proper subcommand-based architecture
- Add `rm-comments go` subcommand for Go comment removal
- Implement global flags: `--recursive/-r`, `--dry-run/-n`, `--verbose/-v`, `--include` (repeatable), `--exclude` (repeatable), `--remove-directives`
- Support stdin/stdout mode via `-` path argument
- Use `go/ast` and `go/parser` for AST-aware comment removal (both `//` and `/* */`)
- Use `go/format` (gofmt) for output formatting
- Preserve Go compiler directives by default using pattern matching (`//go:*` and `// +build`)
- Skip files with no comments (preserve timestamps)
- Implement proper exit codes (0=success, 1=error, 2=misuse)

## Impact
- Affected specs: `cli` (new), `go-comment-removal` (new)
- Affected code: `main.go` (complete rewrite), new files in root package (flat structure)
- Test fixtures: `testdata/` directory
