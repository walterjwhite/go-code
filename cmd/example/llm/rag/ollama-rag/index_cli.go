package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/rag/indexing"
)

func runIndex(args []string) {
	fs, cfg := newCommonFlagSet("index")
	var path string
	var namespace string
	var gitRepo string
	var branch string
	var tag string
	var logLimit int
	var maxFileBytes int

	fs.StringVar(&path, "path", "", "local directory path to index")
	fs.StringVar(&namespace, "namespace", "", "namespace override")
	fs.StringVar(&gitRepo, "git-repo", "", "git repository path to index")
	fs.StringVar(&branch, "branch", "", "git branch name")
	fs.StringVar(&tag, "tag", "", "git tag name")
	fs.IntVar(&logLimit, "log-limit", 100, "number of git commits to index as log docs")
	fs.IntVar(&maxFileBytes, "max-file-bytes", 200000, "max file size in bytes to index from git")
	mustParse(fs, args)

	if path == "" && gitRepo == "" {
		logging.Error(fmt.Errorf("index requires --path or --git-repo"))
	}
	if path != "" && gitRepo != "" {
		logging.Error(fmt.Errorf("choose only one source: --path or --git-repo"))
	}
	if branch != "" && tag != "" {
		logging.Error(fmt.Errorf("use only one of --branch or --tag"))
	}

	rt := newRuntime(cfg, false)
	ctx := context.Background()

	var docs []schema.Document
	if path != "" {
		docs = indexing.IndexLocalPath(ctx, path, namespace)
	} else {
		docs = indexing.IndexGitRepository(gitRepo, branch, tag, namespace, logLimit, maxFileBytes)
	}

	if len(docs) > 0 {
		log.Info().Msgf("sample document metadata: %+v", docs[0].Metadata)
	}

	actualNamespace := namespace
	if actualNamespace == "" && len(docs) > 0 {
		if ns, ok := docs[0].Metadata["namespace"].(string); ok && ns != "" {
			actualNamespace = ns
			log.Info().Msgf("using namespace from document metadata: %s", actualNamespace)
		}
	}

	if actualNamespace == "" {
		log.Warn().Msg("no namespace specified and none found in document metadata, using default collection")
	}

	log.Info().Msgf("indexing %d documents to namespace: '%s'", len(docs), actualNamespace)

	store, err := ensureCollection(ctx, cfg, rt.embedder, actualNamespace)
	if err != nil {
		log.Error().Err(err).Msgf("failed to ensure collection for namespace '%s'", actualNamespace)
		logging.Error(fmt.Errorf("ensure collection: %w", err))
	}

	log.Info().Msgf("successfully created/connected to collection for namespace: '%s'", actualNamespace)
	addDocuments(ctx, store, docs)
}
