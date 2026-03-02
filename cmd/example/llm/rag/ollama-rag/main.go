package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/walterjwhite/go-code/lib/application"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnPanic()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	switch cmd {
	case "index":
		runIndex(args)
	case "search":
		runSearch(args)
	case "ask", "query":
		runAsk(args)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func newCommonFlagSet(name string) (*flag.FlagSet, *commonFlags) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	cfg := &commonFlags{}
	fs.StringVar(&cfg.qdrantURL, "qdrant-url", "http://localhost:6333", "qdrant url")
	fs.StringVar(&cfg.qdrantCollection, "qdrant-collection", "default", "qdrant collection")
	fs.StringVar(&cfg.ollamaURL, "ollama-url", "http://localhost:11434", "ollama url")
	fs.StringVar(&cfg.embedModel, "embed-model", "nomic-embed-text:latest", "embedding model name")
	fs.StringVar(&cfg.model, "model", "mistral:latest", "llm model name")
	return fs, cfg
}

type commonFlags struct {
	qdrantURL        string
	qdrantCollection string
	ollamaURL        string
	embedModel       string
	model            string
}
