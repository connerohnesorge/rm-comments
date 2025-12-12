#!/usr/bin/env node

import * as ts from 'typescript';
import * as fs from 'fs';

function removeComments(filename: string): string {
  const content = fs.readFileSync(filename, 'utf8');
  
  // Determine script kind based on extension for JSX handling
  let scriptKind = ts.ScriptKind.TS;
  if (filename.endsWith('.tsx')) {
    scriptKind = ts.ScriptKind.TSX;
  } else if (filename.endsWith('.jsx')) {
    scriptKind = ts.ScriptKind.JSX;
  } else if (filename.endsWith('.js')) {
    scriptKind = ts.ScriptKind.JS;
  }

  // Remove comments by transpiling with removeComments option
  const result = ts.transpileModule(content, {
    compilerOptions: {
      removeComments: true,
      target: ts.ScriptTarget.Latest,
      module: ts.ModuleKind.ESNext,
      jsx: scriptKind === ts.ScriptKind.TSX || scriptKind === ts.ScriptKind.JSX 
        ? ts.JsxEmit.Preserve 
        : undefined,
    },
    fileName: filename,
  });

  return result.outputText;
}

function main() {
  if (process.argv.length < 3) {
    console.error('Usage: node index.js <file>');
    process.exit(1);
  }

  const filename = process.argv[2];
  
  try {
    const result = removeComments(filename);
    process.stdout.write(result);
  } catch (error) {
    console.error('Error processing file:', error);
    process.exit(1);
  }
}

main();
