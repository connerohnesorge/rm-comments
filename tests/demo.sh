#!/bin/bash

# Demonstration script showing before/after for each language

echo "========================================"
echo "rm-comments Demo - Before & After"
echo "========================================"
echo

for lang in go ts py rs; do
    file="examples/sample.$lang"
    
    case $lang in
        go) lang_name="Go" ;;
        ts) lang_name="TypeScript" ;;
        py) lang_name="Python" ;;
        rs) lang_name="Rust" ;;
    esac
    
    echo "--- $lang_name ($file) ---"
    echo
    echo "BEFORE:"
    echo "--------"
    cat "$file"
    echo
    echo "AFTER:"
    echo "--------"
    ./rm-comments "$file"
    echo
    echo
done
