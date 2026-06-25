package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

func probeVectorSize(ctx context.Context, embedder embeddings.Embedder) (int, error) {
	vecs, err := embedder.EmbedQuery(ctx, "probe")
	if err != nil {
		return 0, err
	}
	if len(vecs) == 0 {
		return 0, fmt.Errorf("embedder returned an empty vector")
	}
	return len(vecs), nil
}

func initStore(u string, collection string, embedder embeddings.Embedder) (qdrant.Store, error) {
	qURL, err := url.Parse(u)
	if err != nil {
		return qdrant.Store{}, fmt.Errorf("parse qdrant URL: %w", err)
	}
	return qdrant.New(
		qdrant.WithURL(*qURL),
		qdrant.WithCollectionName(collection),
		qdrant.WithEmbedder(embedder),
	)
}

func ensureCollection(ctx context.Context, qdrantURL, name string, vectorSize int) error {
	if name == "" {
		return nil
	}
	httpClient := &http.Client{Timeout: 15 * time.Second}
	collURL := fmt.Sprintf("%s/collections/%s", qdrantURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, collURL, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot reach Qdrant at %s: %w", qdrantURL, err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Debug().Msgf("collection %q already exists", name)
		return nil
	}

	body := fmt.Sprintf(`{"vectors":{"size":%d,"distance":"Cosine"}}`, vectorSize)
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, collURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("create collection returned HTTP %d: %v", resp.StatusCode, result)
	}

	log.Info().Msgf("created collection %q (%d dimensions)", name, vectorSize)
	return nil
}
