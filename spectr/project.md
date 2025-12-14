# rm-comments Context

## Purpose
A multi-language comment removal CLI that uses native AST/parsers for each supported language. Distributed exclusively via Nix, which provides both the main CLI and language-specific parser runtimes. Each language is processed via its own subcommand, ensuring accurate comment removal using the language's native parsing capabilities.

## Tech Stack
- **Go 1.25+** - Main CLI binary and built-in Go comment removal (using `go/ast`, `go/parser`)
- **Kong** - CLI argument parsing framework (with shell completion generation)
- **Nix** - Packaging, distribution, and language runtime management
- **treefmt-nix** - Multi-language code formatting
- **golangci-lint** - Go linting

### Language Parser Implementations
Each language parser is a Nix-wrapped binary invoked as a subprocess by the main CLI:

| Language | Subcommand | Parser | Extensions |
|----------|------------|--------|------------|
| Go | `go` | Built-in (`go/ast`, `go/parser`) | `.go` |
| Python | `python` | Python `tokenize` module | `.py`, `.pyw` |
| Rust | `rust` | State-machine based | `.rs` |
| TypeScript | `ts` | TypeScript compiler API | `.ts`, `.tsx` |
| JavaScript | `js` | TypeScript compiler API | `.js`, `.jsx`, `.mjs`, `.cjs` |
| C++ | `cpp` | libclang | `.cpp`, `.hpp`, `.cc`, `.hh`, `.cxx`, `.hxx` |
| C | `c` | libclang | `.c`, `.h` |
| Java | `java` | JavaParser | `.java` |
| Nix | `nix` | rnix-parser | `.nix` |

## CLI Interface

### Command Structure
```
rm-comments <subcommand> [flags] <path>...
rm-comments go [flags] <path>...
rm-comments python [flags] <path>...
# etc.
```

### Global Flags
| Flag | Short | Description |
|------|-------|-------------|
| `--recursive` | `-r` | Recurse into directories |
| `--dry-run` | `-n` | Preview changes without modifying (shows diff) |
| `--verbose` | `-v` | Show per-file status |
| `--include` | | Glob pattern to include (e.g., `--include='*.go'`) |
| `--exclude` | | Glob pattern to exclude (e.g., `--exclude='*_test.go'`) |
| `--remove-directives` | | Remove compiler directives (//go:build, #pragma, etc.) |
| `--remove-shebangs` | | Remove shebang lines (Python only, preserved by default) |
| `-` | | Read from stdin, write to stdout |

### Default Excluded Directories (with -r)
- `.git`
- `vendor`
- `node_modules`
- `.direnv`
- `__pycache__`
- `target` (Rust build dir)

### Exit Codes
- `0` - Success
- `1` - General error (parse error, file not found, etc.)
- `2` - Misuse of command (invalid flags, missing arguments)

### Output Behavior
- **Default**: Silent on success, errors to stderr
- **With `-v`**: Print each modified file path
- **With `--dry-run`**: Show diff-style preview of changes
- **Parse errors**: Fail immediately, do not continue processing

## Project Conventions

### Code Style
- Go: `gofmt`, `golines`, `goimports` (via treefmt-nix)
- Nix: `alejandra` formatter
- Keep code simple and focused - no over-engineering
- Avoid unnecessary abstractions; prefer clear, direct implementations

### Architecture Patterns
- **Subcommand-based language selection**: `rm-comments go <path>`
- **Explicit recursion**: Use `-r`/`--recursive` flag for directory traversal
- **In-place modification by default**: Files are modified directly; use `--dry-run` to preview
- **No backups**: Rely on git for recovery
- **Preserve directives by default**: Compiler directives kept unless `--remove-directives` passed
- **Remove doc comments**: Javadoc, Rustdoc, Python docstrings are treated as regular comments and removed
- **Fail-fast error handling**: Stop on first parse error, exit with error code
- **Stdin/stdout support**: Use `-` as path to read from stdin and write to stdout

### Configuration
- **Optional config file**: `~/.config/rm-comments/config.yaml` or `.rm-comments.yaml` in project root
- **CLI flags override config**: Command-line arguments take precedence
- **Format**: YAML

### Testing Strategy
- **Unit tests**: Go table-driven tests for core parsing/stripping logic
- **Integration tests**: Shell-based tests running actual CLI against fixture files
- **Golden files**: Expected output stored as reference files for comparison
- **Nix checks**: `nix flake check` runs all tests
- **All tests must pass before merge**: Tests must be verified to work (per CLAUDE.md)

### Git Workflow
- **Conventional commits**: `feat:`, `fix:`, `chore:`, `docs:` prefixes
- **Main branch**: `main` (protected)
- **Feature branches**: Branch from main, PR back to main

## Domain Context
- Comment removal must be AST-aware to avoid stripping string literals that look like comments
- Each language has different comment syntaxes and edge cases:
  - **Go**: `//`, `/* */`, plus build directives (`//go:build`, `// +build`, `//go:generate`, etc.)
  - **Python**: `#` comments, docstrings (triple-quoted), shebangs (`#!`)
  - **Rust**: `//`, `/* */` (nested!), doc comments (`///`, `//!`, `/** */`)
  - **TypeScript/JavaScript**: `//`, `/* */`, JSX comments `{/* */}`
  - **C/C++**: `//`, `/* */`, preprocessor directives (`#pragma`, etc.)
  - **Java**: `//`, `/* */`, Javadoc `/** */`
  - **Nix**: `#`, `/* */`
- Preserving compiler directives is critical for build correctness (unless explicitly removed)
- Shebangs in Python are preserved by default for script executability

## Important Constraints
- **Nix-only distribution**: Not published to crates.io, npm, or pkg.go.dev
- **One language per invocation**: Each subcommand handles exactly one language
- **No mixed-language file support**: Users must run appropriate subcommand for each file type
- **Preserve file permissions**: Maintain original file mode when rewriting
- **Handle encoding**: Assume UTF-8, preserve BOM if present
- **Fail on parse errors**: Do not attempt "best effort" on malformed files

## Nix Flake Structure

### Outputs
- `packages.default` - Main `rm-comments` CLI with all languages
- `packages.rm-comments-python` - Standalone Python comment remover
- `packages.rm-comments-rust` - Standalone Rust comment remover
- `packages.rm-comments-ts` - Standalone TypeScript/JavaScript comment remover
- `packages.rm-comments-cpp` - Standalone C++ comment remover
- `packages.rm-comments-c` - Standalone C comment remover
- `packages.rm-comments-java` - Standalone Java comment remover
- `packages.rm-comments-nix` - Standalone Nix comment remover
- `checks` - Test suite run via `nix flake check`
- `devShells.default` - Development environment

### Installation
```bash
# Run directly
nix run github:connerohnesorge/rm-comments -- go -r .

# Install to profile
nix profile install github:connerohnesorge/rm-comments

# Add to flake inputs
inputs.rm-comments.url = "github:connerohnesorge/rm-comments";
```

## External Dependencies
- **Nix/Nixpkgs**: Runtime management and packaging
- **Language runtimes** (provided via Nix):
  - Python 3.x for Python parser
  - Node.js + TypeScript for TS/JS parser
  - Rust toolchain for Rust parser (compile-time only)
  - libclang for C/C++ parser
  - JDK + JavaParser for Java parser
  - rnix-parser (Rust) for Nix parser
- **No external services**: Fully offline/local operation
