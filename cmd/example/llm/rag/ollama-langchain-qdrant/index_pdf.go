package main

import (
	"strings"

	"github.com/gen2brain/go-fitz"

	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func indexPDF(fileName string) ([]schema.Document, error) {
	pagesList := make([]schema.Document, 0)

	doc, err := fitz.New(fileName)
	if err != nil {
		return nil, err
	}

	defer closeResource(doc)

	for idx := range doc.NumPage() {
		text, _ := doc.Text(idx)

		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ToLower(text)

		newDoc := schema.Document{
			PageContent: text,
			Metadata:    map[string]interface{}{"file_path": fileName},
		}

		pagesList = append(pagesList, newDoc)
	}

	return pagesList, nil
}

func closeResource(doc *fitz.Document) {
	logging.Warn(doc.Close(), "doc.Close")
}
