package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func RetrieveDocuments(ctx context.Context, store vectorstores.VectorStore, queryText string, docs int, threshold float32, namespace string) []schema.Document {
	optionsVector := []vectorstores.Option{vectorstores.WithScoreThreshold(threshold)}
	if namespace != "" {
		optionsVector = append(optionsVector, vectorstores.WithFilters(qdrantNamespaceFilter(namespace)))
	}

	retriever := vectorstores.ToRetriever(store, docs, optionsVector...)
	docsRetrieved, err := retriever.GetRelevantDocuments(ctx, queryText)
	if err != nil {
		logging.Error(err)
	}

	log.Info().Msgf("%d docs retrieved", len(docsRetrieved))
	return docsRetrieved
}

func BuildRAGPrompt(question string, docs []schema.Document) string {
	var b strings.Builder
	b.WriteString("Answer the question using only the provided context. ")
	b.WriteString("If the answer is not in the context, say you do not know.\n\n")
	b.WriteString("Context:\n")

	for i, doc := range docs {
		fmt.Fprintf(&b, "[%d] namespace=%s source=%s\n", i+1, MetadataString(doc.Metadata, "namespace"), sourceLabel(doc.Metadata))
		b.WriteString(doc.PageContent)
		b.WriteString("\n\n")
	}

	b.WriteString("Question: ")
	b.WriteString(question)
	b.WriteString("\nAnswer:")
	return b.String()
}

func sourceLabel(metadata map[string]any) string {
	if v := MetadataString(metadata, "file_path"); v != "" {
		return v
	}
	if v := MetadataString(metadata, "commit_hash"); v != "" {
		return v
	}
	return MetadataString(metadata, "doc_type")
}

func MetadataString(metadata map[string]any, key string) string {
	if metadata == nil {
		return ""
	}
	v, ok := metadata[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprint(v)
	}
	return s
}

func Snippet(text string, max int) string {
	text = strings.TrimSpace(text)
	if len(text) <= max {
		return text
	}
	return text[:max] + "..."
}

func qdrantNamespaceFilter(namespace string) map[string]any {
	return map[string]any{
		"must": []map[string]any{
			{
				"key": "namespace",
				"match": map[string]any{
					"value": namespace,
				},
			},
		},
	}
}
