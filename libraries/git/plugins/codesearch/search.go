package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/git"
)

func Search(ctx context.Context, w *git.WorkTreeConfig, pattern string) {
	search := new(w)
	search.NewSearch(ctx, pattern).Search()
}
