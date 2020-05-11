package task

import (
	"context"
)

func (t *Task) Cancel(ctx context.Context) {
	t.changeStatus(ctx, Cancelled)
}
