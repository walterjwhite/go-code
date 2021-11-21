package main

import (
	"context"

	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/vanguard"
)

var (
	session = &vanguard.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	session.Login(context.Background())
}
