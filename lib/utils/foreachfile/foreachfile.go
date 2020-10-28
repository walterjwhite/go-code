package foreachfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
)

type FileMatcher interface {
	MatchesDirectory(root, filePath string, f os.FileInfo) bool
	Matches(root, filePath string, f os.FileInfo) bool

	ToString() string
}

type FilenamePatternMatcher struct {
	Patterns []string
}

type HiddenFileExcluder struct {
	Patterns []string
}

func (m *FilenamePatternMatcher) Matches(root, filePath string, f os.FileInfo) bool {
	return matchesPattern(filePath, m.Patterns...)
}

func (m *FilenamePatternMatcher) MatchesDirectory(root, filePath string, f os.FileInfo) bool {
	return true
}

func (m *FilenamePatternMatcher) ToString() string {
	return fmt.Sprintf("%v", m.Patterns)
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

func (e *HiddenFileExcluder) Matches(root, filePath string, f os.FileInfo) bool {
	return strings.Index(f.Name(), ".") != 0
}

func (e *HiddenFileExcluder) MatchesDirectory(root, filePath string, f os.FileInfo) bool {
	return strings.Index(f.Name(), ".") != 0
}

func (e *HiddenFileExcluder) ToString() string {
	return fmt.Sprintf("%v", e.Patterns)
}

func Execute(root string, callback func(filePath string), patterns ...string) {
	ExecuteCallback(root, callback, &FilenamePatternMatcher{Patterns: patterns})
}

func ExecuteCallback(root string, callback func(filePath string), fileMatcher FileMatcher) {
	files, err := ioutil.ReadDir(root)
	logging.Panic(err)

	for _, f := range files {
		filePath := filepath.Join(root, f.Name())

		if f.IsDir() {
			if fileMatcher.MatchesDirectory(root, filePath, f) {
				ExecuteCallback(filePath, callback, fileMatcher)
			}
		} else {
			if fileMatcher.Matches(root, filePath, f) {
				doFile(filePath, callback, fileMatcher)
			}
		}
	}
}

func doFile(filePath string, callback func(filePath string), fileMatcher FileMatcher) {
	log.Debug().Msgf("calling callback on: %v :%v", filePath, fileMatcher.ToString)
	callback(filePath)
}
