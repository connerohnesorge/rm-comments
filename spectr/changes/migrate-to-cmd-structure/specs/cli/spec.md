## ADDED Requirements

### Requirement: cmd Package Organization
The CLI SHALL be organized with a `cmd/` package containing `root.go` for shared configuration (global flags, exit codes, CLI struct) and `<subcmd>.go` for each language subcommand.

#### Scenario: Package structure organization
- WHEN a new language subcommand is added
- THEN the subcommand is defined in `cmd/<language>.go`
- AND the comment removal logic is implemented in `internal/<language>/`

#### Scenario: Go subcommand file location
- WHEN inspecting the Go subcommand implementation
- THEN the subcommand definition exists in `cmd/go.go`
- AND the comment removal logic exists in `internal/golang/`

### Requirement: internal Package for Language Logic
The CLI SHALL place language-specific comment removal logic in `internal/<language>/` packages, keeping subcommand definitions in `cmd/` minimal and focused on CLI concerns.

#### Scenario: Go comment logic separation
- WHEN the Go subcommand processes a file
- THEN it delegates to `internal/golang` for comment removal
- AND the `cmd/go.go` file contains only CLI-specific code (argument handling, file traversal)
