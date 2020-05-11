package codesearch

import (
	"github.com/walterjwhite/go-application/libraries/codesearch"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"path/filepath"
)

func getIndex(t *task.Task) *codesearch.Instance {
	return &codesearch.Instance{ContentPath: []string{t.Path}, IndexPath: filepath.Join(t.Path, ".codesearch")}
}
