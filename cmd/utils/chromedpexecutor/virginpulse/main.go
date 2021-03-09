package main

import (
	"github.com/walterjwhite/go/lib/application"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/virginpulse"
)

var (
	session = &virginpulse.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	session.Authenticate(application.Context)
	session.RunScript()
	session.Logout()
}
