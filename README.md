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

### Development

```bash
nix develop
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
