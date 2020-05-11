package git

import (
	"context"

	"github.com/walterjwhite/go-application/libraries/logging"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (c *WorkTreeConfig) Checkout(parentCtx context.Context, branch string) {
	logging.Panic(c.W.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(branch)}))
}
