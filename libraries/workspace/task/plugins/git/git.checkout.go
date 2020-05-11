package git

import (
	"context"
	gitl "github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
)

func (g *Git) Checkout(ctx context.Context, t *task.Task, name string) {
	arguments := make(map[string]string)

	// TODO: should be configured via plugins
	arguments["targetDirectory"] = filepath.Join(t.Path, "git", name)
	arguments["branch"] = g.SourceBranch

	gitl.Checkout(ctx, g.Url, arguments)

	// checkout new branch to work on
	gitl.CheckoutBranch(ctx, g.SourceBranch)
}
