package indexing

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func IndexLocalPath(ctx context.Context, rootPath string, namespace string) []schema.Document {
	absPath, err := filepath.Abs(rootPath)
	logging.Error(err)

	if namespace == "" {
		namespace = absPath
	}

	loader := documentloaders.NewRecursiveDirLoader(documentloaders.WithRoot(absPath))
	docs, err := loader.Load(ctx)
	if err != nil {
		logging.Error(fmt.Errorf("load local docs: %w", err))
	}

	for i := range docs {
		if docs[i].Metadata == nil {
			docs[i].Metadata = map[string]any{}
		}
		docs[i].Metadata["source"] = "files"
		docs[i].Metadata["namespace"] = namespace
		docs[i].Metadata["root_path"] = absPath
	}

	log.Info().Msgf("indexed %d local docs in namespace %s", len(docs), namespace)
	return docs
}
