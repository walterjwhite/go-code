package task

import (
	"context"
)

func (t *Task) Complete(ctx context.Context) {
	t.changeStatus(ctx, Completed)
}
