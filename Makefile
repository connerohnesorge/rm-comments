.PHONY: all build clean test backends typescript rust go

all: backends build

build: go

go:
	go build -o rm-comments

backends: typescript rust

typescript:
	cd backends/typescript && npm install && npm run build

rust:
	cd backends/rust && cargo build --release
	mkdir -p bin
	cp backends/rust/target/release/rust-rm-comments bin/

test: all
	./tests/test_all.sh

clean:
	rm -f rm-comments
	rm -rf bin/
	rm -rf backends/typescript/dist/
	rm -rf backends/typescript/node_modules/
	rm -rf backends/rust/target/

install: all
	mkdir -p $(DESTDIR)/usr/local/bin
	mkdir -p $(DESTDIR)/usr/local/libexec/rm-comments
	cp rm-comments $(DESTDIR)/usr/local/bin/
	cp backends/python/python_rm_comments.py $(DESTDIR)/usr/local/libexec/rm-comments/
	cp -r backends/typescript/dist $(DESTDIR)/usr/local/libexec/rm-comments/typescript
	cp bin/rust-rm-comments $(DESTDIR)/usr/local/libexec/rm-comments/

help:
	@echo "Available targets:"
	@echo "  all        - Build all backends and main CLI (default)"
	@echo "  build      - Build main Go CLI"
	@echo "  backends   - Build all language backends"
	@echo "  typescript - Build TypeScript backend"
	@echo "  rust       - Build Rust backend"
	@echo "  test       - Run integration tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install to system (use DESTDIR to override prefix)"
	@echo "  help       - Show this help message"
