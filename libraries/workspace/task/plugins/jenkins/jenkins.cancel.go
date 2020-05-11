package jenkins

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
)

func Cancel(ctx context.Context, t *task.Task, name string) {
	getJenkinsJob(t, name).Cancel(ctx)
}
