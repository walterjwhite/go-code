package run

import (
	"context"
	runl "github.com/walterjwhite/go-application/libraries/run"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"github.com/walterjwhite/go-application/libraries/workspace/task/plugins"
)

type config struct {
	task *task.Task
}

func (c *config) Load(a *runl.Application, prefix string) {
	plugins.Configure(c.task, prefix, a)
}

func Run(ctx context.Context, t *task.Task, name ...string) {
	runl.Configurer = &config{task: t}

	i := runl.New(name...)

	// TODO: parameterize this
	// (local,dev,test,qa,prod ...)
	profile := "default"

	i.Run(ctx, profile)
}
