package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func ValidateFilePath(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if strings.Contains(absPath, "..") {
		return fmt.Errorf("path traversal detected")
	}

	return nil
}

type QAConfig struct {
	Questions []QAItem `yaml:"questions"`
}

type QAItem struct {
	Pattern string `yaml:"pattern"`
	Answer  string `yaml:"answer"`
}

func LoadQAConfig(filePath string) (*QAConfig, error) {
	if filePath == "" {
		return &QAConfig{Questions: []QAItem{}}, nil
	}

	if err := ValidateFilePath(filePath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config QAConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *QAConfig) ConvertToMap() map[string]string {
	qnaMap := make(map[string]string)
	for _, item := range c.Questions {
		qnaMap[item.Pattern] = item.Answer
	}
	return qnaMap
}

func CreateDefaultConfig(filePath string) error {
	if err := ValidateFilePath(filePath); err != nil {
		return err
	}

	defaultConfig := QAConfig{
		Questions: []QAItem{
			{
				Pattern: "authorized",
				Answer:  "Yes",
			},
			{
				Pattern: "background check",
				Answer:  "Yes",
			},
			{
				Pattern: "salary expectation",
				Answer:  "Based on market rates and my experience",
			},
			{
				Pattern: "why do you want to work here",
				Answer:  "I am excited about the opportunity to contribute to your team and grow professionally.",
			},
			{
				Pattern: "experience",
				Answer:  "I have extensive experience in this field.",
			},
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}
