package task

import (
	"context"
)

func Pending(ctx context.Context, submodulePath string) *Task {
	return changeStatus(ctx, submodulePath, "active", "pending")
}
