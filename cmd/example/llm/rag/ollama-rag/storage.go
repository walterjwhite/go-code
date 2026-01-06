package main

import (
	"context"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"net/url"
)

var (
	qdrantUrlFlag        = flag.String("qdrant-url", "http://localhost:6333", "qdrant url")
	qdrantCollectionFlag = flag.String("qdrant-collection", "default", "qdrant collection")

	store *qdrant.Store
)

func storage() {
	qdrantUrl, err := url.Parse(*qdrantUrlFlag)
	logging.Panic(err)

	qdrant, err := qdrant.New(
		qdrant.WithURL(*qdrantUrl),
		qdrant.WithCollectionName(*qdrantCollectionFlag),
		qdrant.WithEmbedder(embedder),
	)

	logging.Panic(err)

	store = &qdrant
}

func index(dirFile string) {
	log.Info().Msg(*documentPathFlag)

	l := documentloaders.NewRecursiveDirLoader(documentloaders.WithRoot(dirFile) /*, documentloaders.WithAllowExts(".go", ".md", ".txt", ".java", ".csv", ".sh")*/)

	docs, err := l.Load(context.Background())
	logging.Warn(err, "loadDocuments")

	add(docs)
}

func add(docs []schema.Document) {
	if len(docs) == 0 {
		return
	}

	log.Info().Msgf("adding %v docs", len(docs))
	_, err := store.AddDocuments(context.Background(), docs)
	logging.Warn(err, "useStorage-AddDocuments")
}
