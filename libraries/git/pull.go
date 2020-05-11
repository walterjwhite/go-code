package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"

	"gopkg.in/src-d/go-git.v4"
)

func (c *WorkTreeConfig) Pull(parentCtx context.Context) {
	d := 30 * time.Second
	timeout.Limit(c.doPull, &d, parentCtx)
}

// only supports fast-forwards
func (c *WorkTreeConfig) doPull() {
	err := c.W.Pull(&git.PullOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		logging.Panic(err)
	}
}
