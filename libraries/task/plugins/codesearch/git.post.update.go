package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/task"
)

func (c *Codesearch) PostUpdate(ctx context.Context, t *task.Task) {
	c.Index(ctx, t)
}
