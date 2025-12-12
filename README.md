# rm-comments

A multi-language comment stripping utility packaged using Nix.

## Philosophy

The tool embodies a "right tool for the job" philosophy. Rather than building one universal parser that attempts to understand every language's syntax (inevitably getting edge cases wrong), each backend uses native or first-party parsing infrastructure for its target language:

- **Go**: Uses Go's `go/parser` and `go/ast` packages
- **TypeScript/JavaScript**: Uses the TypeScript Compiler API
- **Python**: Uses Python's built-in `tokenize` module
- **Rust**: Uses `syn` and `proc-macro2` crates

This approach ensures accurate comment removal that respects each language's unique syntax and edge cases.

## Supported Languages

- Go (`.go`)
- TypeScript (`.ts`, `.tsx`)
- JavaScript (`.js`, `.jsx`)
- Python (`.py`)
- Rust (`.rs`)

## Installation

### Using Nix Flakes

```bash
nix build
nix run . -- <file>
```

### Using Make

```bash
# Build all backends and main CLI
make all

# Run tests
make test

# Install to system
sudo make install

# Clean build artifacts
make clean
```

### Development

With Nix:
```bash
nix develop
```

Without Nix:
```bash
# Ensure you have Go, Node.js, Python 3, and Rust/Cargo installed
make all
```

## Usage

```bash
rm-comments <file>
```

The tool automatically detects the language based on file extension and uses the appropriate native parser to remove comments.

### Examples

```bash
# Remove comments from a Go file
rm-comments examples/sample.go

# Remove comments from a TypeScript file
rm-comments examples/sample.ts

# Remove comments from a Python file
rm-comments examples/sample.py

# Remove comments from a Rust file
rm-comments examples/sample.rs
```

## Architecture

The project consists of:

1. **Main CLI** (`main.go`): Orchestrates the different language backends
2. **Go Backend** (`go_backend.go`): Native Go parser integration
3. **TypeScript Backend** (`backends/typescript/`): TypeScript Compiler API integration
4. **Python Backend** (`backends/python/`): Python tokenizer integration
5. **Rust Backend** (`backends/rust/`): Syn parser integration

Each backend is implemented using the language's own parsing tools to ensure accuracy and proper handling of edge cases.

## Technical Details

### Go Backend
Uses Go's `go/parser` package to parse source files into an Abstract Syntax Tree (AST), then uses `go/printer` to regenerate the code without comments. By not including the `parser.ParseComments` flag, comments are automatically excluded from the AST.

### TypeScript/JavaScript Backend
Uses the TypeScript Compiler API's `transpileModule` function with the `removeComments: true` option. This leverages TypeScript's sophisticated understanding of JavaScript/TypeScript syntax to safely remove all comment types while preserving code semantics.

### Python Backend
Uses Python's built-in `tokenize` module, which is the same tokenizer used by the Python interpreter. It filters out `COMMENT` tokens while preserving all other tokens, ensuring perfect compatibility with Python's syntax rules.

### Rust Backend
Uses the `syn` crate, Rust's de-facto standard parsing library (used by procedural macros). The parsed AST is converted back to tokens using `quote::ToTokens`, which automatically excludes comments. The output is then formatted with `rustfmt` for readability.

## Building Backends

### TypeScript Backend

```bash
cd backends/typescript
npm install
npm run build
```

### Rust Backend

```bash
cd backends/rust
cargo build --release
```

### Main Go CLI

```bash
go build -o rm-comments
```

## License

MIT
