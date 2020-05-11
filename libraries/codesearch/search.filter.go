package codesearch

import (
	"regexp"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type SearchFileProcessor interface {
	Matches(name string) bool
}

type SearchFileInclusionFilter struct {
	Patterns []string
}

func (s *SearchInstance) filterFiles(post []uint32) []uint32 {
	fnames := make([]uint32, 0, len(post))

	for _, fileid := range post {
		name := s.ix.Name(fileid)
		if s.SearchFileProcessor.Matches(name) {
			fnames = append(fnames, fileid)
		}
	}

	return fnames
}

func (f *SearchFileInclusionFilter) Matches(name string) bool {
	for _, pattern := range f.Patterns {
		matched, err := regexp.MatchString(pattern, name)
		logging.Panic(err)

		if matched {
			return true
		}
	}

	return false
}
