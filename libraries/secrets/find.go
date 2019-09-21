package secrets

import (
	"flag"
	"io/ioutil"
	"log"

	"strings"

	"github.com/walterjwhite/go-application/libraries/logging"
)

type Search struct {
	IncludeDeprecated bool
	Patterns          []string
}

type NoSearchCriteriaError struct{}

func (e *NoSearchCriteriaError) Error() string {
	return "You must specify at least one pattern to search."
}

func Find(patterns []string, callback func(filePath string)) {
	doFind(SecretsConfigurationInstance.RepositoryPath, patterns, callback)
}

func doFind(root string, patterns []string, callback func(filePath string)) {
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

	log.Printf("searching in: %v\n", SecretsConfigurationInstance.RepositoryPath)
	log.Printf("patterns: %v\n", patterns)

	if len(patterns) == 0 {
		logging.Panic(&NoSearchCriteriaError{})
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
