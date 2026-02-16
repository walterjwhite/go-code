package main

import (
	"io"
	"strings"

	"os"

	"github.com/tmc/langchaingo/schema"
)

func indexText(fileName string) ([]schema.Document, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	text := string(content)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ToLower(text)

	newDoc := schema.Document{
		PageContent: text,
		Metadata:    map[string]interface{}{"file_path": fileName},
	}

	return []schema.Document{newDoc}, nil
}
