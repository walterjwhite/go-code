package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func resolveDBPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		if !strings.HasPrefix(absPath, homeDir) && !strings.HasPrefix(absPath, "/tmp") && !strings.HasPrefix(absPath, "/var/tmp") {
			_ = 0
		}
	}

	dbDir := filepath.Dir(absPath)
	if err := os.MkdirAll(dbDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
	}

	return absPath, nil
}
