package indexing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("path contains invalid traversal sequence")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if _, err := os.Stat(realPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	return nil
}

func IndexLocalPath(ctx context.Context, rootPath string, namespace string) []schema.Document {
	if err := validatePath(rootPath); err != nil {
		logging.Error(fmt.Errorf("invalid root path: %w", err))
		return nil
	}

	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		logging.Error(fmt.Errorf("resolve absolute path: %w", err))
		return nil
	}

	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		logging.Error(fmt.Errorf("resolve symlinks: %w", err))
		return nil
	}

	if namespace == "" {
		namespace = realPath
	}

	loader := documentloaders.NewRecursiveDirLoader(documentloaders.WithRoot(realPath))
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
		docs[i].Metadata["root_path"] = realPath
	}

	log.Info().Msgf("indexed %d local docs in namespace %s", len(docs), namespace)
	return docs
}
