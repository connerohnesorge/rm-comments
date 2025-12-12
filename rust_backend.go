package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func processRustFile(filename string, content []byte) (string, error) {
	// Find the Rust backend
	backendPath := findRustBackend()
	if backendPath == "" {
		return "", fmt.Errorf("Rust backend not found")
	}

	// Execute the Rust backend
	cmd := exec.Command(backendPath, filename)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("Rust backend failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute Rust backend: %w", err)
	}

	return string(output), nil
}

func findRustBackend() string {
	// Check if running from Nix installation
	execPath, err := os.Executable()
	if err == nil {
		nixPath := filepath.Join(filepath.Dir(execPath), "..", "libexec", "rm-comments", "rust-rm-comments")
		if _, err := os.Stat(nixPath); err == nil {
			return nixPath
		}
	}

	// Check local development path
	localPath := filepath.Join("backends", "rust", "target", "release", "rust-rm-comments")
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// Check bin directory
	binPath := filepath.Join("bin", "rust-rm-comments")
	if _, err := os.Stat(binPath); err == nil {
		return binPath
	}

	return ""
}
