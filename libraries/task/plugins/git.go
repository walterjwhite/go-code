package plugins

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/task"
)

// TODO: go-git doesn't support hooks, instead, simplify this
type PreCommit interface {
	PreCommit(ctx context.Context, t *task.Task)
}

type PostCommit interface {
	PostCommit(ctx context.Context, t *task.Task)
}
