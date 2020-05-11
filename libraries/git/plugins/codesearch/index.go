package codesearch

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/git"
)

func Index(ctx context.Context, w *git.WorkTreeConfig) {
	search := new(w)
	search.NewDefaultIndex(ctx).Index()
}
