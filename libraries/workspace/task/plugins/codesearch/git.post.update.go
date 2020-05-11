package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
)

func PostUpdate(ctx context.Context, t *task.Task) {
	Index(ctx, t)
}
