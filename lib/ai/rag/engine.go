package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/textsplitter"
)

func (cfg *Config) New(ctx context.Context) error {
	embedLLM, err := ollama.New(ollama.WithModel(cfg.EmbedModel), ollama.WithServerURL(cfg.OllamaURL))
	if err != nil {
		return fmt.Errorf("create embed LLM: %w", err)
	}
	embedder, err := embeddings.NewEmbedder(embedLLM)
	if err != nil {
		return fmt.Errorf("create embedder: %w", err)
	}
	vectorSize, err := probeVectorSize(ctx, embedder)
	if err != nil {
		return fmt.Errorf("probe embedding dimension: %w", err)
	}

	if err := ensureCollection(ctx, cfg.QdrantURL, cfg.Collection, vectorSize); err != nil {
		return fmt.Errorf("ensure main collection: %w", err)
	}
	if err := ensureCollection(ctx, cfg.QdrantURL, cfg.CacheCollection, vectorSize); err != nil {
		return fmt.Errorf("ensure cache collection: %w", err)
	}

	store, err := initStore(cfg.QdrantURL, cfg.Collection, embedder)
	if err != nil {
		return err
	}
	cacheStore, err := initStore(cfg.QdrantURL, cfg.CacheCollection, embedder)
	if err != nil {
		return err
	}
	chatLLM, err := ollama.New(ollama.WithModel(cfg.ChatModel), ollama.WithServerURL(cfg.OllamaURL))
	if err != nil {
		return fmt.Errorf("create chat LLM: %w", err)
	}

	cfg.Engine = &Engine{
		cfg:        cfg,
		store:      store,
		cacheStore: cacheStore,
		llm:        chatLLM,
		splitter: textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(cfg.ChunkSize),
			textsplitter.WithChunkOverlap(cfg.ChunkOverlap),
			textsplitter.WithSeparators([]string{"\n\n", "\n", ". ", " ", ""}),
		),
	}

	return nil
}
