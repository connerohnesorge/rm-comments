package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func processPythonFile(filename string, content []byte) (string, error) {
	// Find the Python backend
	backendPath := findPythonBackend()
	if backendPath == "" {
		return "", fmt.Errorf("Python backend not found")
	}

	// Execute the Python backend
	cmd := exec.Command("python3", backendPath, filename)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("Python backend failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute Python backend: %w", err)
	}

	return string(output), nil
}

func findPythonBackend() string {
	// Check if running from Nix installation
	execPath, err := os.Executable()
	if err == nil {
		nixPath := filepath.Join(filepath.Dir(execPath), "..", "libexec", "rm-comments", "python_rm_comments.py")
		if _, err := os.Stat(nixPath); err == nil {
			return nixPath
		}
	}

	// Check local development path
	localPath := filepath.Join("backends", "python", "python_rm_comments.py")
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	return ""
}
