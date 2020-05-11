package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"

	"gopkg.in/src-d/go-git.v4"
)

func (c *WorkTreeConfig) Push(parentCtx context.Context) {
	d := 30 * time.Second
	timeout.Limit(c.doPush, &d, parentCtx)
}

func (c *WorkTreeConfig) doPush() {
	logging.Panic(c.R.Push(&git.PushOptions{}))
}
