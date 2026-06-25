package rag

import (
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

type Config struct {
	QdrantURL       string
	OllamaURL       string
	Collection      string
	CacheCollection string // Added for semantic caching collection isolation
	EmbedModel      string
	ChatModel       string
	ChunkSize       int
	ChunkOverlap    int
	IngestTimeout   time.Duration
	NumDocs         int
	CacheThreshold  float32 // e.g., 0.93

	Engine *Engine
}

type Engine struct {
	cfg        *Config
	store      vectorstores.VectorStore
	cacheStore vectorstores.VectorStore // Added for semantic caching
	llm        llms.Model
	splitter   textsplitter.TextSplitter
}
