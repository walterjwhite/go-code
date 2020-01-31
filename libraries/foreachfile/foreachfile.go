package foreachfile

import (
	"io/ioutil"
	"path/filepath"

	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func Execute(root string, callback func(filePath string), patterns ...string) {
	doFind(root, callback, patterns...)
}

func doFind(root string, callback func(filePath string), patterns ...string) {
	files, err := ioutil.ReadDir(root)
	logging.Panic(err)

	for _, f := range files {
		filePath := filepath.Join(root, f.Name())

		if f.IsDir() {
			doFind(filePath, callback, patterns...)
		} else {
			doFile(filePath, callback, patterns...)
		}
	}
}

func doFile(filePath string, callback func(filePath string), patterns ...string) {
	/*
	   if !strings.HasSuffix(filePath, "/value") {
	           return
	   }
	*/

	if !matchesPattern(filePath, patterns...) {
		return
	}

	log.Debug().Msgf("calling callback on: %v :%v", filePath, patterns)
	callback(filePath)
}

func matchesPattern(filePath string, patterns ...string) bool {
	if len(patterns) == 0 {
		return true
	}

	for _, pattern := range patterns {
		if !strings.Contains(filePath, pattern) {
			return false
		}
	}

	return true
}
