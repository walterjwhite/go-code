package secrets

import (
	"errors"
	"flag"
	"github.com/rs/zerolog/log"
	"io/ioutil"

	"strings"

	"github.com/walterjwhite/go-application/libraries/logging"
)

func Find(patterns []string, callback func(filePath string)) {
	doFind(SecretsConfigurationInstance.RepositoryPath, patterns, callback)
}

func doFind(root string, patterns []string, callback func(filePath string)) {
	initialize()

	files, err := ioutil.ReadDir(root)
	logging.Panic(err)

	for _, f := range files {
		filePath := root + "/" + f.Name()

		if f.IsDir() {
			doFind(filePath, patterns, callback)
		} else {
			findOne(filePath, patterns, callback)
		}
	}
}

func NewFind() []string {
	// this should NOT be needed
	//flag.Parse()

	patterns := flag.Args()

	log.Debug().Msgf("searching in: %v", SecretsConfigurationInstance.RepositoryPath)
	log.Debug().Msgf("patterns: %v", patterns)

	if len(patterns) == 0 {
		logging.Panic(errors.New("You must specify at least one pattern to search."))
	}

	return patterns
}

func findOne(filePath string, patterns []string, callback func(filePath string)) {
	if !strings.HasSuffix(filePath, "/value") {
		return
	}

	for _, pattern := range patterns {
		if !strings.Contains(filePath, pattern) {
			return
		}
	}

	callback(filePath)
}
