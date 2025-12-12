use std::env;
use std::fs;
use std::io::Write;
use std::process::{Command, Stdio};
use quote::ToTokens;

fn remove_comments(content: &str) -> Result<String, Box<dyn std::error::Error>> {
    // Parse the Rust file using syn
    let syntax_tree = syn::parse_file(content)?;
    
    // Convert back to tokens, which excludes comments
    // The ToTokens trait implementation automatically strips comments
    let tokens = syntax_tree.to_token_stream();
    
    // Convert tokens back to string
    let code = tokens.to_string();
    
    // Try to format with rustfmt if available
    match format_with_rustfmt(&code) {
        Ok(formatted) => Ok(formatted),
        Err(_) => Ok(code), // Fall back to unformatted if rustfmt fails
    }
}

fn format_with_rustfmt(code: &str) -> Result<String, Box<dyn std::error::Error>> {
    let mut child = Command::new("rustfmt")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()?;
    
    if let Some(mut stdin) = child.stdin.take() {
        stdin.write_all(code.as_bytes())?;
    }
    
    let output = child.wait_with_output()?;
    
    if output.status.success() {
        Ok(String::from_utf8(output.stdout)?)
    } else {
        Err("rustfmt failed".into())
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();
    
    if args.len() < 2 {
        eprintln!("Usage: {} <file>", args[0]);
        std::process::exit(1);
    }
    
    let filename = &args[1];
    
    // Read the file
    let content = match fs::read_to_string(filename) {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Error reading file: {}", e);
            std::process::exit(1);
        }
    };
    
    // Remove comments
    match remove_comments(&content) {
        Ok(result) => print!("{}", result),
        Err(e) => {
            eprintln!("Error processing file: {}", e);
            std::process::exit(1);
        }
    }
}
