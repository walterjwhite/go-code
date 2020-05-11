package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
)

// TODO:
// 1. support toggling case-sensitivity
// 2. support file pattern
func Search(ctx context.Context, t *task.Task, pattern string) {
	getIndex(t).NewSearch(ctx, pattern).Search()
}
