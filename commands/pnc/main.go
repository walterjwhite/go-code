package main

import (
	"errors"
	"flag"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"github.com/walterjwhite/go-application/libraries/utils/web/chromedpexecutor/plugins/pnc"
)

var (
	session = &pnc.Session{}
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
		log.Debug().Msgf("username: %v", session.Credentials.Username)
		session.Login(application.Context)
	case "logout":
		session.Logout()
	}

	time.Sleep(10 * time.Minute)
}
