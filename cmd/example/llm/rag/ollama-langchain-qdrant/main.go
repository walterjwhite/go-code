package main

import (
	"flag"

	"net/url"
	"os"

	"github.com/tmc/langchaingo/embeddings"

	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	qdrantUrlFlag = flag.String("qdrantUrl", "http://localhost:6333", "qdrant url")
	ollamaUrlFlag = flag.String("ollamaUrl", "http://127.0.0.1:11434", "ollama url")
	modelFlag     = flag.String("model", "mistral", "model to use [mistral]")


	e      *embeddings.EmbedderImpl
	urlAPI *url.URL
	llm *ollama.LLM
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnPanic()
	var err error
	llm, err = NewOllama()
	logging.Error(err)

	e, err = embeddings.NewEmbedder(llm)
	logging.Error(err)

	urlAPI, err = url.Parse(*qdrantUrlFlag)
	logging.Error(err)

	cmd := cliCommand()

	logging.Error(cmd.Run(application.Context, os.Args))
}
