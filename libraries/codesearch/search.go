package codesearch

import (
	"context"
	"github.com/google/codesearch/index"
	"github.com/google/codesearch/regexp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

type SearchInstance struct {
	Context         context.Context
	Pattern         string
	CaseInsensitive bool

	SearchFileProcessor   SearchFileProcessor
	SearchOutputProcessor SearchOutputProcessor

	Instance *Instance

	regexpPattern string
	regexp        *regexp.Regexp
	ix            *index.Index

	buf []byte
}

func (i *Instance) NewSearch(ctx context.Context, pattern string) *SearchInstance {
	return &SearchInstance{Context: ctx, Pattern: pattern, Instance: i, SearchOutputProcessor: &DefaultSearchOutputProcessor{Output: os.Stdout}}
}

func (s *SearchInstance) Search() {
	log.Info().Msgf("Searching %v", s.Instance.IndexPath)

	s.updatePattern()

	s.ix = index.Open(s.Instance.IndexPath)

	re, err := regexp.Compile(s.regexpPattern)
	logging.Panic(err)

	s.regexp = re

	q := index.RegexpQuery(re.Syntax)
	s.searchFiles(s.setupSearch(q))
}

func (s *SearchInstance) updatePattern() {
	s.regexpPattern = "(?m)" + s.Pattern
	if s.CaseInsensitive {
		s.regexpPattern = "(?i)" + s.regexpPattern
	}
}

func (s *SearchInstance) setupSearch(q *index.Query) []uint32 {
	post := s.ix.PostingQuery(q)

	if s.SearchFileProcessor == nil {
		return post
	}

	return s.filterFiles(post)
}

func (s *SearchInstance) searchFiles(post []uint32) {
	for _, fileid := range post {
		name := s.ix.Name(fileid)
		s.file(name)
	}
}
