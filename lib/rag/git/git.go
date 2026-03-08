package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return fmt.Errorf("repository path does not exist: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("repository path is not a directory")
	}

	return nil
}

func validateGitRef(ref string) error {
	if strings.TrimSpace(ref) == "" {
		return fmt.Errorf("git reference cannot be empty")
	}

	validRef := regexp.MustCompile(`^[a-zA-Z0-9/_\-.]+$`)
	if !validRef.MatchString(ref) {
		return fmt.Errorf("git reference contains invalid characters")
	}

	if strings.Contains(ref, "..") {
		return fmt.Errorf("git reference cannot contain '..'")
	}

	dangerousChars := []string{";", "|", "&", "$", "`", "(", ")", "{", "}", "<", ">", "!", "\\", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(ref, char) {
			return fmt.Errorf("git reference contains dangerous character")
		}
	}

	return nil
}

func validateFilePath(file string) error {
	if strings.TrimSpace(file) == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if strings.Contains(file, "..") {
		return fmt.Errorf("file path cannot contain '..'")
	}

	if filepath.IsAbs(file) {
		return fmt.Errorf("file path must be relative")
	}

	dangerousChars := []string{";", "|", "&", "$", "`", "(", ")", "{", "}", "<", ">", "!", "\\", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(file, char) {
			return fmt.Errorf("file path contains dangerous character")
		}
	}

	return nil
}

func ShowFile(repoPath, ref, file string) ([]byte, error) {
	if err := validateGitRef(ref); err != nil {
		return nil, fmt.Errorf("invalid git reference: %w", err)
	}
	if err := validateFilePath(file); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

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

	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repository path: %w", err)
	}

	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repository path: %w", err)
	}

	for _, arg := range args {
		if err := validateGitArg(arg); err != nil {
			return nil, err
		}
	}

	cmd := exec.Command("git", append([]string{"-C", realPath}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

func validateGitArg(arg string) error {
	if strings.TrimSpace(arg) == "" {
		return fmt.Errorf("git argument cannot be empty")
	}

	dangerousChars := []string{";", "|", "&", "$", "`", "(", ")", "{", "}", "<", ">", "!", "\\", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(arg, char) {
			return fmt.Errorf("git argument contains invalid character")
		}
	}

	return nil
}

func IsBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) >= 0
}
