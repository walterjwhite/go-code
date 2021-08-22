package main

import (
	"context"

	"github.com/walterjwhite/go/lib/application"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/patientnotebook"
)

var (
	session = &patientnotebook.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	session.Login(context.Background())
}
