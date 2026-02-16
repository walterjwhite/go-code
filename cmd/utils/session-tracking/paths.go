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

	dbDir := filepath.Dir(path)
	if err := os.MkdirAll(dbDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create database directory %s: %w", dbDir, err)
	}

	return path, nil
}
