package secrets

import (
	"errors"
	"flag"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-application/libraries/foreachfile"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func Find(callback func(filePath string), patterns ...string) {
	doFind(SecretsConfigurationInstance.RepositoryPath, callback, patterns...)
}

func doFind(root string, callback func(filePath string), patterns ...string) {
	initialize()

	foreachfile.Execute(root, callback, patterns...)
}

func NewFind() []string {
	patterns := flag.Args()

	if len(patterns) == 0 {
		logging.Panic(errors.New("You must specify at least one pattern to search."))
	}

	initialize()

	log.Debug().Msgf("searching in: %v", SecretsConfigurationInstance.RepositoryPath)
	log.Debug().Msgf("patterns: %v", patterns)

	patterns = append(patterns, "/value")

	return patterns
}
