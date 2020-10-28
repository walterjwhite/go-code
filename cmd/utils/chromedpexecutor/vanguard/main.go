package main

import (
	"context"
	"errors"
	"flag"

	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"

	"github.com/walterjwhite/go/lib/utils/web/chromedpexecutor/plugins/vanguard"
)

var (
	session = &vanguard.Session{}
)

func init() {
	application.ConfigureWithProperties(session)
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) < 1 {
		logging.Panic(errors.New("Command is required (login, logout)"))
	}

	switch flag.Args()[0] {
	case "login":
		session.Login(context.Background())
	}
}
