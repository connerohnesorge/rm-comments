## ADDED Requirements

### Requirement: Subcommand-Based CLI Structure
The CLI SHALL use a subcommand-based structure where each supported language has its own subcommand (e.g., `rm-comments go`, `rm-comments python`).

#### Scenario: User invokes Go subcommand
- WHEN user runs `rm-comments go <path>`
- THEN the Go comment removal logic processes the specified path

#### Scenario: User invokes without subcommand
- WHEN user runs `rm-comments` without a subcommand
- THEN help text is displayed showing available subcommands

### Requirement: Global Recursive Flag
The CLI SHALL support a `--recursive` (short: `-r`) flag that enables directory traversal.

#### Scenario: Recursive flag processes directories
- WHEN user runs `rm-comments go -r ./src`
- THEN all `.go` files in `./src` and its subdirectories are processed

#### Scenario: Without recursive flag on directory
- WHEN user runs `rm-comments go ./src` without `-r`
- THEN an error is returned indicating a directory was provided without recursive flag

### Requirement: Global Dry-Run Flag
The CLI SHALL support a `--dry-run` (short: `-n`) flag that previews changes without modifying files.

#### Scenario: Dry-run shows diff
- WHEN user runs `rm-comments go -n file.go`
- THEN a diff-style preview of changes is displayed to stdout
- AND the original file is not modified

### Requirement: Global Verbose Flag
The CLI SHALL support a `--verbose` (short: `-v`) flag that prints the path of each modified file to stdout (path only, no status prefix).

#### Scenario: Verbose mode shows modified files
- WHEN user runs `rm-comments go -v ./src -r`
- AND files `src/main.go` and `src/util.go` are modified
- THEN stdout contains:
  ```
  src/main.go
  src/util.go
  ```

#### Scenario: Verbose mode skips unchanged files
- WHEN user runs `rm-comments go -v ./src -r`
- AND file `src/empty.go` has no comments
- THEN `src/empty.go` is not printed to stdout

### Requirement: Include Pattern Flag
The CLI SHALL support a repeatable `--include` flag accepting glob patterns to filter which files are processed. Multiple `--include` flags SHALL be combined with OR logic.

#### Scenario: Include pattern filters files
- WHEN user runs `rm-comments go -r --include='*_test.go' ./src`
- THEN only files matching `*_test.go` are processed

#### Scenario: Multiple include patterns
- WHEN user runs `rm-comments go -r --include='*_test.go' --include='*_bench.go' ./src`
- THEN files matching either `*_test.go` OR `*_bench.go` are processed

### Requirement: Exclude Pattern Flag
The CLI SHALL support a repeatable `--exclude` flag accepting glob patterns to exclude files from processing. Multiple `--exclude` flags SHALL be combined with OR logic.

#### Scenario: Exclude pattern skips files
- WHEN user runs `rm-comments go -r --exclude='*_test.go' ./src`
- THEN files matching `*_test.go` are skipped

#### Scenario: Multiple exclude patterns
- WHEN user runs `rm-comments go -r --exclude='*_test.go' --exclude='*_bench.go' ./src`
- THEN files matching either `*_test.go` OR `*_bench.go` are skipped

### Requirement: Remove Directives Flag
The CLI SHALL support a `--remove-directives` flag that removes compiler directives that would otherwise be preserved.

#### Scenario: Remove directives flag
- WHEN user runs `rm-comments go --remove-directives file.go`
- AND file.go contains `//go:build linux`
- THEN the directive is removed along with regular comments

### Requirement: Default Excluded Directories
The CLI SHALL exclude common non-source directories by default when using recursive mode: `.git`, `vendor`, `node_modules`, `.direnv`, `__pycache__`, `target`.

#### Scenario: Default exclusions applied
- WHEN user runs `rm-comments go -r ./project`
- AND `./project` contains a `vendor/` directory
- THEN files in `vendor/` are not processed

### Requirement: Stdin/Stdout Mode
The CLI SHALL support reading from stdin and writing to stdout when `-` is provided as the path argument.

#### Scenario: Process from stdin to stdout
- WHEN user runs `echo "// comment\nfunc main() {}" | rm-comments go -`
- THEN the output is written to stdout with comments removed
- AND no files are modified on disk

### Requirement: Exit Codes
The CLI SHALL exit with code 0 on success, 1 on general errors (parse error, file not found), and 2 on misuse (invalid flags, missing arguments).

#### Scenario: Successful execution
- WHEN user runs `rm-comments go file.go` and processing succeeds
- THEN exit code is 0

#### Scenario: Parse error
- WHEN user runs `rm-comments go invalid.go` and the file has syntax errors
- THEN exit code is 1
- AND error message is written to stderr

#### Scenario: Invalid flags
- WHEN user runs `rm-comments go --invalid-flag`
- THEN exit code is 2
- AND usage help is displayed
