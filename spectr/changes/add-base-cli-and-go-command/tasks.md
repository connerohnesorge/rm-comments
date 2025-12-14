## 1. CLI Framework Setup
- [ ] 1.1 Replace placeholder CLI structure with proper Kong subcommand architecture
- [ ] 1.2 Define global flags struct (recursive, dry-run, verbose, include, exclude, remove-directives)
- [ ] 1.3 Define Go subcommand struct with path arguments
- [ ] 1.4 Implement exit code handling (0=success, 1=error, 2=misuse)

## 2. Go Comment Removal Core
- [ ] 2.1 Create Go comment removal package using `go/ast` and `go/parser`
- [ ] 2.2 Implement single-line comment (`//`) removal
- [ ] 2.3 Implement multi-line comment (`/* */`) removal
- [ ] 2.4 Implement directive detection and preservation logic (//go:build, //go:generate, // +build, etc.)
- [ ] 2.5 Preserve file formatting (indentation, spacing) after comment removal

## 3. File Processing
- [ ] 3.1 Implement single file processing
- [ ] 3.2 Implement recursive directory traversal with `-r` flag
- [ ] 3.3 Implement include/exclude glob pattern filtering
- [ ] 3.4 Implement default directory exclusions (.git, vendor, node_modules, .direnv)
- [ ] 3.5 Implement stdin/stdout mode with `-` path argument
- [ ] 3.6 Preserve file permissions when rewriting

## 4. Output Modes
- [ ] 4.1 Implement silent mode (default behavior)
- [ ] 4.2 Implement verbose mode (`-v`) showing modified file paths
- [ ] 4.3 Implement dry-run mode (`-n`) showing diff preview

## 5. Error Handling
- [ ] 5.1 Implement fail-fast on parse errors
- [ ] 5.2 Add proper error messages to stderr
- [ ] 5.3 Handle missing files, permission errors gracefully

## 6. Testing
- [ ] 6.1 Create test fixtures in `testdata/` directory with various Go comment patterns
- [ ] 6.2 Write table-driven unit tests for comment removal logic
- [ ] 6.3 Write integration tests for CLI behavior
- [ ] 6.4 Test directive preservation scenarios (//go:* pattern matching)
- [ ] 6.5 Test stdin/stdout mode
- [ ] 6.6 Test unchanged file skip behavior
