package task

import (
	"context"
)

func (t *Task) Active(ctx context.Context, submodulePath string) *Task {
	return changeStatus(ctx, submodulePath, "canceled", "active")
}
