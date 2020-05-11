package codesearch

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
	"regexp"
)

type IndexFileProcessor interface {
	IsInclude(directory, basename string, info os.FileInfo) bool
}

type FileExclusionFilter struct {
	Patterns []string
}

func (f *FileExclusionFilter) IsInclude(directory, basename string, info os.FileInfo) bool {
	for _, pattern := range f.Patterns {
		matched, err := regexp.MatchString(pattern, basename)
		logging.Panic(err)

		if matched {
			log.Debug().Msgf("Excluding: %v via pattern (%v)", basename, pattern)
			return false
		}
	}

	log.Debug().Msgf("Including: %v", basename)
	return true
}
