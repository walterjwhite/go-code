package rag

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (e *Engine) Query(ctx context.Context, question string, out io.Writer) ([]schema.Document, error) {
	if cached, err := e.checkCache(ctx, question); err == nil && cached != "" {
		_, werr := io.WriteString(out, cached)
		return nil, werr
	}

	docs, err := e.store.SimilaritySearch(ctx, question, e.cfg.NumDocs)
	if err != nil || len(docs) == 0 {
		_, werr := fmt.Fprintln(out, "(no relevant documents found - try ingesting more content)")
		return nil, werr
	}

	var buf strings.Builder
	_, err = e.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{llms.TextParts(llms.ChatMessageTypeHuman, buildPrompt(question, docs))},
		llms.WithStreamingFunc(func(_ context.Context, chunk []byte) error {
			if _, werr := out.Write(chunk); werr != nil {
				return werr
			}
			buf.Write(chunk)
			return nil
		}),
	)
	if err == nil && e.cfg.CacheCollection != "" {
		_, err = e.saveCache(ctx, question, buf.String())
		logging.Warn(err, "saveCache")
	}
	return docs, err
}

func (e *Engine) QueryText(ctx context.Context, question string) (string, []schema.Document, error) {
	var buf bytes.Buffer
	docs, err := e.Query(ctx, question, &buf)
	if err != nil {
		return "", docs, err
	}
	return strings.TrimSpace(buf.String()), docs, nil
}

func buildPrompt(question string, docs []schema.Document) string {
	var sb strings.Builder
	sb.WriteString("You are a helpful assistant. Answer the question using ONLY the context below.\n")
	sb.WriteString("If the context does not contain enough information, say so - do not invent facts.\n\n### Context\n\n")
	for i, doc := range docs {
		src, _ := doc.Metadata["source"].(string)
		fmt.Fprintf(&sb, "--- Excerpt %d (source: %s) ---\n%s\n\n", i+1, src, doc.PageContent)
	}
	fmt.Fprintf(&sb, "### Question\n\n%s\n\n### Answer\n\n", question)
	return sb.String()
}
