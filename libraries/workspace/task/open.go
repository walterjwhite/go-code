package task

import (
	"context"
)

func (t *Task) Open(ctx context.Context) {
	t.changeStatus(ctx, Open)
}
