#!/usr/bin/env python3
"""
Python comment removal backend using Python's native tokenize module.
This uses the first-party tokenizer to accurately parse Python source code.
"""

import sys
import tokenize
import io


def remove_comments(filename):
    """Remove comments from Python source file using tokenize module."""
    with open(filename, 'rb') as f:
        tokens = tokenize.tokenize(f.readline)
        
        result = []
        prev_toktype = tokenize.INDENT
        prev_end = (1, 0)
        
        for toktype, tokval, start, end, line in tokens:
            # Skip comment and encoding tokens
            if toktype in (tokenize.COMMENT, tokenize.ENCODING):
                continue
            
            # Handle newlines
            if toktype == tokenize.NEWLINE or toktype == tokenize.NL:
                result.append(tokval)
                prev_end = end
                continue
            
            # Add proper spacing for indentation
            if start[0] > prev_end[0]:  # New line
                if toktype == tokenize.INDENT:
                    result.append(tokval)
                else:
                    result.append(' ' * start[1])
            elif start[1] > prev_end[1]:  # Same line, add spaces
                result.append(' ' * (start[1] - prev_end[1]))
            
            # Add the token value (skip INDENT tokens as we handle them above)
            if toktype != tokenize.INDENT:
                result.append(tokval)
            
            prev_toktype = toktype
            prev_end = end
    
    return ''.join(result)


def main():
    if len(sys.argv) < 2:
        print('Usage: python3 python_rm_comments.py <file>', file=sys.stderr)
        sys.exit(1)
    
    filename = sys.argv[1]
    
    try:
        result = remove_comments(filename)
        sys.stdout.write(result)
    except Exception as e:
        print(f'Error processing file: {e}', file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()
