package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/rag/query"
)

func runSearch(args []string) {
	fs, cfg := newCommonFlagSet("search")
	var queryText string
	var namespace string
	var docs int
	var threshold float64

	fs.StringVar(&queryText, "query", "", "search query")
	fs.StringVar(&namespace, "namespace", "", "optional namespace")
	fs.IntVar(&docs, "docs", 5, "# of docs to include")
	fs.Float64Var(&threshold, "threshold", 0.4, "score threshold [0..1]")
	mustParse(fs, args)

	if queryText == "" {
		logging.Error(fmt.Errorf("search requires --query"))
	}

	rt := newRuntime(cfg, false)

	store := rt.store
	var namespaceFilter string
	if namespace != "" {
		var err error
		store, err = ensureCollection(context.Background(), cfg, rt.embedder, namespace)
		if err != nil {
			logging.Error(fmt.Errorf("ensure collection: %w", err))
		}
		namespaceFilter = ""
	} else {
		namespaceFilter = namespace
	}

	found := query.RetrieveDocuments(context.Background(), store, queryText, docs, float32(threshold), namespaceFilter)

	if len(found) == 0 {
		fmt.Println("no documents matched")
		return
	}

	for i, doc := range found {
		source := query.MetadataString(doc.Metadata, "file_path")
		if source == "" {
			source = query.MetadataString(doc.Metadata, "doc_type")
		}
		fmt.Printf("[%d] score=%.4f source=%s namespace=%s\n", i+1, doc.Score, source, query.MetadataString(doc.Metadata, "namespace"))
		fmt.Println(query.Snippet(doc.PageContent, 280))
		fmt.Println()
	}
}

func runAsk(args []string) {
	fs, cfg := newCommonFlagSet("ask")
	var prompt string
	var namespace string
	var docs int
	var threshold float64

	fs.StringVar(&prompt, "prompt", "", "prompt")
	fs.StringVar(&namespace, "namespace", "", "optional namespace")
	fs.IntVar(&docs, "docs", 5, "# of docs to include")
	fs.Float64Var(&threshold, "threshold", 0.4, "score threshold [0..1]")
	mustParse(fs, args)

	if prompt == "" {
		logging.Error(fmt.Errorf("ask requires --prompt"))
	}

	rt := newRuntime(cfg, true)

	store := rt.store
	var namespaceFilter string
	if namespace != "" {
		var err error
		store, err = ensureCollection(context.Background(), cfg, rt.embedder, namespace)
		if err != nil {
			logging.Error(fmt.Errorf("ensure collection: %w", err))
		}
		namespaceFilter = ""
	} else {
		namespaceFilter = namespace
	}

	docsRetrieved := query.RetrieveDocuments(context.Background(), store, prompt, docs, float32(threshold), namespaceFilter)
	if len(docsRetrieved) == 0 {
		log.Warn().Msg("no documents retrieved")
		fmt.Println("no supporting documents matched")
		return
	}

	answerPrompt := query.BuildRAGPrompt(prompt, docsRetrieved)
	answer, err := llms.GenerateFromSinglePrompt(context.Background(), rt.llm, answerPrompt)
	if err != nil {
		logging.Error(err)
	}

	fmt.Println(answer)
}

func mustParse(fs *flag.FlagSet, args []string) {
	err := fs.Parse(args)
	if err != nil {
		logging.Error(err)
	}
}
