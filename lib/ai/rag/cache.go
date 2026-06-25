package rag

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

const threshold = float32(0.97)

func (e *Engine) checkCache(ctx context.Context, q string) (string, error) {
	if e.cfg.CacheCollection == "" {
		return "", nil
	}

	res, err := e.cacheStore.SimilaritySearch(ctx, q, 1, vectorstores.WithScoreThreshold(threshold))
	if err != nil || len(res) == 0 {
		return "", err
	}

	cachedQuestion, ok := res[0].Metadata["original_question"].(string)
	if !ok {
		return "", nil // Fallback to live LLM if metadata is corrupted
	}

	words := strings.Fields(strings.ToLower(q))
	if len(words) > 2 {
		subjectWord := words[len(words)-2]
		if !strings.Contains(strings.ToLower(cachedQuestion), subjectWord) {
			log.Debug().Msgf("Semantic cache hit rejected due to entity mismatch: %q vs %q", q, cachedQuestion)
			return "", nil // Force a cache miss!
		}
	}

	return res[0].Metadata["answer"].(string), nil
}

func (e *Engine) saveCache(ctx context.Context, q, ans string) ([]string, error) {
	return e.cacheStore.AddDocuments(ctx, []schema.Document{{
		PageContent: q, // The prompt acts as the embedded anchor key
		Metadata: map[string]any{
			"answer":            ans,
			"original_question": q,
			"saved_at":          time.Now().Unix(),
		},
	}})
}
