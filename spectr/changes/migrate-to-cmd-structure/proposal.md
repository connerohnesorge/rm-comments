# Change: Migrate CLI to cmd/ Package Structure

## Why
The current implementation has all code in a single `main.go` file (378 lines). As more language subcommands are added (Python, Rust, TypeScript, etc.), this monolithic approach will become difficult to maintain. Migrating to a `cmd/` package structure with `cmd/root.go` for shared configuration and `cmd/<subcmd>.go` for each language subcommand follows Go CLI conventions and improves code organization.

## What Changes
- Extract CLI root configuration and global flags to `cmd/root.go`
- Extract Go subcommand to `cmd/go.go`
- Move comment removal logic to `internal/golang/` package for reuse
- Simplify `main.go` to just parse and execute
- Keep Kong as the CLI framework (no framework change)
- Maintain all existing functionality and behavior

## Impact
- Affected specs: `cli` (MODIFIED - implementation structure)
- Affected code:
  - `main.go` - Simplified to minimal entrypoint
  - `cmd/root.go` - New: CLI root, global flags
  - `cmd/go.go` - New: Go subcommand definition
  - `internal/golang/strip.go` - New: Comment removal logic
  - `internal/golang/diff.go` - New: Diff generation
