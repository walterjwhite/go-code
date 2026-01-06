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
	var err error
	llm, err = NewOllama()
	logging.Panic(err)

	e, err = embeddings.NewEmbedder(llm)
	logging.Panic(err)

	urlAPI, err = url.Parse(*qdrantUrlFlag)
	logging.Panic(err)

	cmd := cliCommand()

	logging.Panic(cmd.Run(application.Context, os.Args))
}
