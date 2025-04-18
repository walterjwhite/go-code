package main

import (
	"context"

	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/authenticate"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider/remote"
)

var (
	s = &authenticate.Session{}
)

func init() {
	application.Configure(s)
}

func main() {
	defer application.OnEnd()

	ctx := context.Background()

	s.With(ctx, remote.New(ctx))

	s.Login()
}
