package git

import (
	"context"

	"github.com/walterjwhite/go-application/libraries/logging"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func (c *WorkTreeConfig) CreateBranch(parentCtx context.Context, sourceBranch, targetBranch string) {
	logging.Panic(c.W.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(sourceBranch), Branch: plumbing.NewBranchReferenceName(targetBranch), Create: true}))
}
