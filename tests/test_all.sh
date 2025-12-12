#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Testing rm-comments multi-language comment stripping utility"
echo "=============================================================="
echo

# Build the project
echo "Building project..."
cd "$(dirname "$0")/.."
go build -o rm-comments
cd backends/typescript && npm install && npm run build && cd ../..
cd backends/rust && cargo build --release && cd ../..
echo

# Test Go backend
echo "Testing Go backend..."
OUTPUT=$(./rm-comments examples/sample.go)
if echo "$OUTPUT" | grep -q "// This is a single line comment"; then
    echo -e "${RED}FAIL: Go comments not removed${NC}"
    exit 1
else
    echo -e "${GREEN}PASS: Go backend${NC}"
fi
echo

# Test TypeScript backend
echo "Testing TypeScript backend..."
OUTPUT=$(./rm-comments examples/sample.ts)
if echo "$OUTPUT" | grep -q "// TypeScript sample file"; then
    echo -e "${RED}FAIL: TypeScript comments not removed${NC}"
    exit 1
else
    echo -e "${GREEN}PASS: TypeScript backend${NC}"
fi
echo

# Test Python backend
echo "Testing Python backend..."
OUTPUT=$(./rm-comments examples/sample.py)
if echo "$OUTPUT" | grep -q "# This is a Python sample file"; then
    echo -e "${RED}FAIL: Python comments not removed${NC}"
    exit 1
else
    echo -e "${GREEN}PASS: Python backend${NC}"
fi
echo

# Test Rust backend
echo "Testing Rust backend..."
OUTPUT=$(./rm-comments examples/sample.rs)
if echo "$OUTPUT" | grep -q "// Rust sample file"; then
    echo -e "${RED}FAIL: Rust comments not removed${NC}"
    exit 1
else
    echo -e "${GREEN}PASS: Rust backend${NC}"
fi
echo

echo -e "${GREEN}All tests passed!${NC}"
