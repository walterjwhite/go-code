package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
	},
}

type runtime struct {
	llm      *ollama.LLM
	embedder *embeddings.EmbedderImpl
	store    *qdrant.Store
}

func newRuntime(cfg *commonFlags, includeLLM bool) *runtime {
	embedLLM, err := ollama.New(
		ollama.WithServerURL(cfg.ollamaURL),
		ollama.WithModel(cfg.embedModel),
	)
	logging.Error(err)

	embedder, err := embeddings.NewEmbedder(embedLLM)
	logging.Error(err)

	qdrantURL, err := url.Parse(cfg.qdrantURL)
	logging.Error(err)

	st, err := qdrant.New(
		qdrant.WithURL(*qdrantURL),
		qdrant.WithCollectionName(cfg.qdrantCollection),
		qdrant.WithEmbedder(embedder),
	)
	logging.Error(err)

	r := &runtime{
		embedder: embedder,
		store:    &st,
	}

	if includeLLM {
		llm, llmErr := ollama.New(
			ollama.WithServerURL(cfg.ollamaURL),
			ollama.WithModel(cfg.model),
		)
		logging.Error(llmErr)
		r.llm = llm
	}

	return r
}

func addDocuments(ctx context.Context, st *qdrant.Store, docs []schema.Document) {
	if len(docs) == 0 {
		log.Warn().Msg("no documents to add")
		return
	}

	log.Info().Msgf("adding %d docs", len(docs))
	ids, err := st.AddDocuments(ctx, docs)
	if err != nil {
		log.Error().
			Err(err).
			Str("error_type", fmt.Sprintf("%T", err)).
			Msgf("failed to add documents: %v", err)
		logging.Error(err)
	}
	log.Info().Msgf("successfully added %d documents with ids: %v", len(ids), ids)
}

func sanitizeCollectionName(namespace string) string {
	if namespace == "" {
		return ""
	}

	parts := strings.Split(namespace, "::")

	var repoName string
	var branch string

	if len(parts) > 0 {
		repoPath := parts[0]
		repoName = filepath.Base(repoPath)
		repoName = strings.TrimSuffix(repoName, ".git")
	}

	if len(parts) > 1 {
		branch = parts[1]
	}

	var collectionName string
	if branch != "" {
		collectionName = fmt.Sprintf("%s.%s", repoName, branch)
	} else {
		collectionName = repoName
	}

	collectionName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			return r
		}
		return '_'
	}, collectionName)

	log.Info().Msgf("sanitized namespace '%s' to collection name '%s'", namespace, collectionName)
	return collectionName
}

func createQdrantCollection(qdrantURL, collectionName string, vectorSize int) error {
	baseURL, err := url.Parse(qdrantURL)
	if err != nil {
		return fmt.Errorf("parse qdrant url: %w", err)
	}
	baseURL.Path = "/collections/" + url.PathEscape(collectionName)
	createURL := baseURL.String()

	payload := map[string]any{
		"vectors": map[string]any{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal collection config: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "PUT", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer close(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusBadRequest {
			log.Info().Msgf("collection '%s' may already exist: %s", collectionName, string(body))
			return nil
		}
		return fmt.Errorf("create collection failed (status %d): %s", resp.StatusCode, string(body))
	}

	log.Info().Msgf("created collection '%s' with vector size %d", collectionName, vectorSize)
	return nil
}

func close(i io.Closer) {
	if i == nil {
		return
	}
	logging.Warn(i.Close(), "close")
}

func ensureCollection(ctx context.Context, cfg *commonFlags, embedder *embeddings.EmbedderImpl, namespace string) (*qdrant.Store, error) {
	collectionName := cfg.qdrantCollection
	if namespace != "" {
		collectionName = sanitizeCollectionName(namespace)
	}

	log.Info().Msgf("ensuring collection: %s (namespace: %s, default: %s)", collectionName, namespace, cfg.qdrantCollection)

	vectorSize := 768
	err := createQdrantCollection(cfg.qdrantURL, collectionName, vectorSize)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to create collection '%s', will try to connect anyway", collectionName)
	}

	qdrantURL, err := url.Parse(cfg.qdrantURL)
	if err != nil {
		return nil, fmt.Errorf("parse qdrant url: %w", err)
	}

	st, err := qdrant.New(
		qdrant.WithURL(*qdrantURL),
		qdrant.WithCollectionName(collectionName),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		return nil, fmt.Errorf("create qdrant store for collection '%s': %w", collectionName, err)
	}

	log.Info().Msgf("connected to collection: %s", collectionName)
	return &st, nil
}

