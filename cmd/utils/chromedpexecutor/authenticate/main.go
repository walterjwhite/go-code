package main

import (
	"context"

	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/authenticate"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"
)

var (
	s = &authenticate.Session{}
)

func init() {
	application.ConfigureWithProperties(s)
}

func main() {
	defer application.OnEnd()

	ctx := context.Background()

	s.With(ctx, remote.New(ctx))

	s.Login()
}
