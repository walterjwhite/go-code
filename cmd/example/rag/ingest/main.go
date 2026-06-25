package main

import (
	"context"
	"flag"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/ai/document"
	"github.com/walterjwhite/go-code/lib/ai/rag"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	ragConfig = &rag.Config{}
)

func init() {
	application.Configure(ragConfig)
	logging.Error(ragConfig.New(application.Context), "init")
}

func main() {
	defer application.OnPanic()

	timeout := ragConfig.IngestTimeout
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(application.Context, timeout)
	defer cancel()

	log.Info().Msg("loading documents…")
	var docs []document.Document
	for _, path := range flag.Args() {
		loaded, err := document.LoadPath(path)
		if err != nil {
			logging.Warn(err, "loading documents")
			continue
		}
		for _, d := range loaded {
			log.Info().Msgf("%s (%s, %d chars)", d.FilePath, d.Format, len(d.Content))
		}
		docs = append(docs, loaded...)
	}
	if len(docs) == 0 {
		log.Warn().Msg("no supported documents found in the given paths")
		return
	}
	log.Info().Msgf("%d document(s) loaded\n", len(docs))

	log.Info().Msg("connected to Qdrant and Ollama")

	log.Info().Msgf("chunking and embedding (chunk-size=%d, overlap=%d, timeout=%s)…\n",
		ragConfig.ChunkSize, ragConfig.ChunkOverlap, timeout)
	n, err := ragConfig.Engine.Ingest(ctx, docs)
	logging.Error(err, "ingest")

	log.Info().Msgf("%d chunk(s) stored in collection %q\n", n, ragConfig.Collection)
}
