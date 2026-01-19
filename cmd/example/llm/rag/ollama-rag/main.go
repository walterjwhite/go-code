package main

import (
	"flag"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	modelFlag        = flag.String("model", "mistral:latest", "model name")
	documentPathFlag = flag.String("doc-path", "", "document path")
	promptFlag       = flag.String("prompt", "", "prompt")

	llm      *ollama.LLM
	embedder *embeddings.EmbedderImpl
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnPanic()
	llm, err := ollama.New(ollama.WithModel(*modelFlag))
	logging.Error(err)


	embedder, err = embeddings.NewEmbedder(llm /*ollamaEmbedderModel*/)
	logging.Error(err)

	storage()

	if len(*documentPathFlag) > 0 {
		index(*documentPathFlag)
	}

	if len(*promptFlag) > 0 {
		query()
	}
}
