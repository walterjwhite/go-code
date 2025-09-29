package main

import (
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/utils/publisher/provider/google"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/linkedin/learning"
)

var (
	session  = &learning.Session{}
	provider *google.Provider
)

func init() {
	application.Configure(session)

	session.Init(application.Context)
	provider = google.New(application.Context)
}

func main() {
	session.Run(provider)
}
