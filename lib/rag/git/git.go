package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func validateRepoPath(repoPath string) error {
	if strings.TrimSpace(repoPath) == "" {
		return fmt.Errorf("repository path cannot be empty")
	}

	if strings.Contains(repoPath, "..") {
		return fmt.Errorf("repository path cannot contain '..'")
	}

	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("repository path does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("repository path is not a directory")
	}

	return nil
}

func ShowFile(repoPath, ref, file string) ([]byte, error) {
	return OutputBytes(repoPath, "show", ref+":"+file)
}

func Lines(repoPath string, args ...string) ([]string, error) {
	out, err := Output(repoPath, args...)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(out, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func Output(repoPath string, args ...string) (string, error) {
	out, err := OutputBytes(repoPath, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func OutputBytes(repoPath string, args ...string) ([]byte, error) {
	if err := validateRepoPath(repoPath); err != nil {
		return nil, err
	}

	absPath, _ := filepath.Abs(repoPath)

	cmd := exec.Command("git", append([]string{"-C", absPath}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

func IsBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) >= 0
}
