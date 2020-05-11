package codesearch

import (
	"context"

	"github.com/google/codesearch/index"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
	"path/filepath"
	"sort"
)

type IndexInstance struct {
	Context            context.Context
	Instance           *Instance
	IndexFileProcessor IndexFileProcessor

	ix *index.IndexWriter
}

func (i *Instance) NewDefaultIndex(ctx context.Context) *IndexInstance {
	return &IndexInstance{Context: ctx, Instance: i, IndexFileProcessor: &FileExclusionFilter{Patterns: []string{"^\\.", "^#", "^~", "~$"}}}
}

func (i *IndexInstance) Index() {
	i.doProcessContentPaths()
	i.ix = index.Create(i.Instance.IndexPath)
	i.ix.AddPaths(i.Instance.ContentPath)

	i.doIndex()

	i.ix.Flush()
}

func (i *IndexInstance) doProcessContentPaths() {
	for j, arg := range i.Instance.ContentPath {
		expanded, err := homedir.Expand(arg)
		logging.Panic(err)

		a, err := filepath.Abs(expanded)
		logging.Panic(err)

		i.Instance.ContentPath[j] = a
	}

	sort.Strings(i.Instance.ContentPath)
	log.Info().Msgf("added paths: %v", i.Instance.ContentPath)
}

func (i *IndexInstance) doIndex() {
	for _, contentPath := range i.Instance.ContentPath {
		logging.Panic(filepath.Walk(contentPath, i.isProcess))
	}
}

func (i *IndexInstance) isProcess(path string, info os.FileInfo, err error) error {
	dir, basename := filepath.Split(path)

	if err != nil {
		log.Warn().Msgf("(err) skipping: %v : %v", path, err)
		return nil
	}

	if !i.IndexFileProcessor.IsInclude(dir, basename, info) {
		if info.IsDir() {
			log.Info().Msgf("isDir: skipping dir: %v", path)
			return filepath.SkipDir
		}

		log.Info().Msgf("skipping file: %v", path)
		return nil
	}

	if info != nil && info.Mode()&os.ModeType == 0 {
		log.Info().Msgf("file: %v", path)

		i.ix.AddFile(path)
	} else {
		log.Info().Msgf("skipping file: %v", path)
	}

	return nil
}
