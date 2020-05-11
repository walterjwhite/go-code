package codesearch

import (
	codesearchl "github.com/walterjwhite/go-application/libraries/codesearch"
	"github.com/walterjwhite/go-application/libraries/git"
	"path/filepath"
)

func new(w *git.WorkTreeConfig) *codesearchl.Instance {
	return &codesearchl.Instance{IndexPath: filepath.Join(w.Path, ".codesearch.index"),
		ContentPath: []string{w.Path}}
}
