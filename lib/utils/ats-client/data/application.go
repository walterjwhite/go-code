package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

func ValidateFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	cleanPath := filepath.Clean(filePath)

	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: relative paths not allowed")
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal detected: resolved path contains parent references")
	}

	if _, err := os.Stat(absPath); err == nil {
		fileInfo, err := os.Lstat(absPath)
		if err != nil {
			return fmt.Errorf("cannot access file: %w", err)
		}
		if !fileInfo.Mode().IsRegular() {
			return fmt.Errorf("path is not a regular file")
		}
	}

	return nil
}

type ApplicationData struct {
	Timestamp   time.Time         `yaml:"timestamp"`
	URL         string            `yaml:"url"`
	AccountInfo AccountInfo       `yaml:"account_info"`
	Questions   map[string]string `yaml:"questions"`
	ResumePath  string            `yaml:"resume_path"`
}

type AccountInfo struct {
	Email     string `yaml:"email"`
	FirstName string `yaml:"first_name"`
	LastName  string `yaml:"last_name"`
	Phone     string `yaml:"phone"`
	Address   string `yaml:"address"`
	City      string `yaml:"city"`
	State     string `yaml:"state"`
	ZipCode   string `yaml:"zip_code"`
	Country   string `yaml:"country"`
}

func SaveApplicationData(data *ApplicationData, filePath string) error {
	if err := ValidateFilePath(filePath); err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, yamlData, 0600)
}

func LoadApplicationData(filePath string) (*ApplicationData, error) {
	if err := ValidateFilePath(filePath); err != nil {
		return nil, err
	}

	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data ApplicationData
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
