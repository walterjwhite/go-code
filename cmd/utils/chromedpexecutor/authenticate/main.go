package main

import (
	"context"

	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/authenticate"
)

var (
	session = &authenticate.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	ctx := context.Background()
	session.Login(ctx)
	session.KeepAlive(ctx)
}
