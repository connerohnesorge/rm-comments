package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func processTypeScriptFile(filename string, content []byte, ext string) (string, error) {
	// Find the TypeScript backend
	backendPath := findTypeScriptBackend()
	if backendPath == "" {
		return "", fmt.Errorf("TypeScript backend not found")
	}

	// Execute the TypeScript backend
	cmd := exec.Command("node", backendPath, filename)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("TypeScript backend failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute TypeScript backend: %w", err)
	}

	return string(output), nil
}

func findTypeScriptBackend() string {
	// Check if running from Nix installation
	execPath, err := os.Executable()
	if err == nil {
		nixPath := filepath.Join(filepath.Dir(execPath), "..", "libexec", "rm-comments", "typescript", "index.js")
		if _, err := os.Stat(nixPath); err == nil {
			return nixPath
		}
	}

	// Check local development path
	localPath := filepath.Join("backends", "typescript", "dist", "index.js")
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	return ""
}
