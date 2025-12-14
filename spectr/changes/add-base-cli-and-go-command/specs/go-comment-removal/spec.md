## ADDED Requirements

### Requirement: AST-Based Comment Removal
The Go comment remover SHALL use Go's native `go/ast` and `go/parser` packages for AST-aware comment removal, ensuring string literals containing comment-like syntax are not affected.

#### Scenario: Comment in string literal preserved
- WHEN processing Go code containing `s := "// not a comment"`
- THEN the string literal is preserved unchanged

#### Scenario: Actual comment removed
- WHEN processing Go code containing `// this is a comment`
- THEN the comment is removed from the output

### Requirement: Single-Line Comment Removal
The Go comment remover SHALL remove single-line comments (`//`).

#### Scenario: Single-line comment removed
- WHEN processing:
  ```go
  func main() {
      // comment to remove
      fmt.Println("hello")
  }
  ```
- THEN output contains:
  ```go
  func main() {
      fmt.Println("hello")
  }
  ```

#### Scenario: End-of-line comment removed
- WHEN processing `x := 1 // inline comment`
- THEN output is `x := 1`

### Requirement: Multi-Line Comment Removal
The Go comment remover SHALL remove multi-line comments (`/* */`).

#### Scenario: Multi-line comment removed
- WHEN processing:
  ```go
  /* This is a
     multi-line comment */
  func main() {}
  ```
- THEN output contains:
  ```go
  func main() {}
  ```

#### Scenario: Inline multi-line comment removed
- WHEN processing `x := /* comment */ 1`
- THEN output is `x := 1`

### Requirement: Directive Preservation by Default
The Go comment remover SHALL preserve Go compiler directives by default using pattern matching: any comment matching `//go:*` pattern or `// +build` legacy format SHALL be preserved.

#### Scenario: Build directive preserved
- WHEN processing Go code containing `//go:build linux`
- AND `--remove-directives` flag is not set
- THEN the directive is preserved in output

#### Scenario: Generate directive preserved
- WHEN processing Go code containing `//go:generate stringer -type=MyType`
- AND `--remove-directives` flag is not set
- THEN the directive is preserved in output

#### Scenario: Legacy build constraint preserved
- WHEN processing Go code containing `// +build linux`
- AND `--remove-directives` flag is not set
- THEN the directive is preserved in output

#### Scenario: Unknown go directive preserved
- WHEN processing Go code containing `//go:newdirective something`
- AND `--remove-directives` flag is not set
- THEN the directive is preserved in output (pattern matching supports future directives)

#### Scenario: Directive removed when flag set
- WHEN processing Go code containing `//go:build linux`
- AND `--remove-directives` flag is set
- THEN the directive is removed

### Requirement: Output Formatting via go/format
The Go comment remover SHALL use `go/format` (gofmt) to format output after comment removal. This ensures consistent, canonical Go formatting.

#### Scenario: Output is gofmt-normalized
- WHEN processing any Go file
- THEN the output is formatted according to gofmt standards

#### Scenario: Indentation normalized
- WHEN processing code with non-standard indentation
- THEN the output uses standard gofmt indentation (tabs)

### Requirement: Doc Comment Removal
The Go comment remover SHALL treat doc comments (comments immediately preceding declarations) as regular comments and remove them.

#### Scenario: Function doc comment removed
- WHEN processing:
  ```go
  // Add returns the sum of a and b.
  func Add(a, b int) int {
      return a + b
  }
  ```
- THEN output contains:
  ```go
  func Add(a, b int) int {
      return a + b
  }
  ```

### Requirement: File Extension Filtering
The Go comment remover SHALL only process files with the `.go` extension.

#### Scenario: Go file processed
- WHEN path `main.go` is provided
- THEN the file is processed

#### Scenario: Non-Go file rejected
- WHEN path `main.rs` is provided to the Go subcommand
- THEN an error is returned indicating wrong file type

### Requirement: In-Place Modification
The Go comment remover SHALL modify files in-place by default (when not using dry-run or stdin mode). Files with no comments to remove SHALL be skipped (not rewritten).

#### Scenario: File modified in-place
- WHEN user runs `rm-comments go file.go`
- AND `file.go` contains comments
- THEN `file.go` is rewritten with comments removed
- AND original file permissions are preserved

#### Scenario: Unchanged file skipped
- WHEN user runs `rm-comments go empty.go`
- AND `empty.go` contains no comments
- THEN `empty.go` is not modified
- AND file timestamp is unchanged

### Requirement: Parse Error Handling
The Go comment remover SHALL fail immediately on parse errors and not attempt best-effort processing.

#### Scenario: Invalid Go syntax
- WHEN processing a file with invalid Go syntax
- THEN an error is returned
- AND no partial output is produced
- AND exit code is 1
