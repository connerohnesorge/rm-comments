# Contributing to rm-comments

Thank you for your interest in contributing to rm-comments!

## Development Setup

### Prerequisites

- Go 1.21 or later
- Node.js 18+ and npm
- Python 3.8+
- Rust 1.70+ and Cargo
- (Optional) Nix with flakes enabled

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/connerohnesorge/rm-comments.git
cd rm-comments
```

2. Build all components:
```bash
make all
```

3. Run tests:
```bash
make test
```

## Project Structure

```
rm-comments/
├── main.go                    # Main CLI entry point
├── go_backend.go              # Go comment removal using go/parser
├── typescript_backend.go      # TypeScript backend dispatcher
├── python_backend.go          # Python backend dispatcher
├── rust_backend.go            # Rust backend dispatcher
├── backends/
│   ├── typescript/            # TypeScript/JavaScript backend
│   │   └── src/index.ts       # TypeScript compiler API implementation
│   ├── python/                # Python backend
│   │   └── python_rm_comments.py  # Python tokenizer implementation
│   └── rust/                  # Rust backend
│       └── src/main.rs        # Syn parser implementation
├── examples/                  # Sample files for each language
└── tests/                     # Integration tests
```

## Adding Support for a New Language

To add support for a new language, follow these steps:

1. **Create a backend**: Add a new directory under `backends/` for your language.

2. **Use native parsing tools**: The philosophy of this project is to use each language's native or first-party parsing infrastructure. Examples:
   - For Java: Use JavaParser or the compiler API
   - For C/C++: Use libclang
   - For Ruby: Use Ripper or parser gem

3. **Create a dispatcher in Go**: Add a new file like `yourlang_backend.go` that:
   - Implements a `processYourLangFile(filename string, content []byte) (string, error)` function
   - Finds and executes your backend binary/script
   - Handles errors appropriately

4. **Update main.go**: Add your file extension to the switch statement in `main.go`.

5. **Add tests**: 
   - Create sample file in `examples/`
   - Add test case to `tests/test_all.sh`

6. **Update documentation**:
   - Add language to README.md
   - Update supported language list
   - Document the backend's implementation

## Code Style

- **Go**: Follow standard Go formatting (`gofmt`)
- **TypeScript**: Use TypeScript's recommended style
- **Python**: Follow PEP 8
- **Rust**: Use `rustfmt`

## Testing

Always run tests before submitting a PR:

```bash
make test
```

Test coverage should include:
- Single-line comments
- Multi-line comments
- Inline comments
- Edge cases specific to the language

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and ensure they pass
5. Update documentation
6. Submit a pull request

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT).
