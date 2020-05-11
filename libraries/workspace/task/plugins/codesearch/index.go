package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
)

func Index(ctx context.Context, t *task.Task) {
	getIndex(t).NewDefaultIndex(ctx).Index()
}
