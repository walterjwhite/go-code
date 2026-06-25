package rag

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/walterjwhite/go-code/lib/ai/document"
)

func (e *Engine) Ingest(ctx context.Context, docs []document.Document) (int, error) {
	var langDocs []schema.Document

	for _, doc := range docs {
		raw := []schema.Document{{
			PageContent: doc.Content,
			Metadata: map[string]any{
				"source":   doc.FilePath,
				"filename": filepath.Base(doc.FilePath),
				"format":   string(doc.Format),
			},
		}}

		chunks, err := textsplitter.SplitDocuments(e.splitter, raw)
		if err != nil {
			return 0, fmt.Errorf("split %s: %w", doc.FilePath, err)
		}

		for i := range chunks {
			chunks[i].Metadata["chunk"] = i
		}

		log.Info().Msgf("loaded: %s: %d chunk(s)", doc.FilePath, len(chunks))
		langDocs = append(langDocs, chunks...)
	}

	if len(langDocs) == 0 {
		return 0, nil
	}

	if _, err := e.store.AddDocuments(ctx, langDocs); err != nil {
		return 0, fmt.Errorf("store documents: %w", err)
	}

	return len(langDocs), nil
}
