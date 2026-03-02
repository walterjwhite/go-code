package main

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/token/providers/cli"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix"
)

var (
	session       = &citrix.Session{}
	firstRun bool = true
)

func init() {
	application.Configure(session)

	log.Debug().Msg("initializing google pubsub")

	application.Load(session.GoogleProvider)
	logging.Error(session.GoogleProvider.Init(application.Context))

	log.Debug().Msgf("conf: %v", session.GoogleProvider)
}

func main() {
	defer application.OnPanic()
	for {
		token := getToken()

		session.Init(application.Context)
		err := session.Run(token)
		if errors.Is(err, context.DeadlineExceeded) {
			log.Warn().Msg("expected context deadline exceeded")
		} else {
			logging.Error(err)
		}

		firstRun = false
	}
}

func getToken() string {
	if firstRun {
		token := cli.New().Get()
		if len(token) > 0 {
			log.Info().Msg("using cmdline token")
			return token
		}
	}

	log.Info().Msg("using google pubsub")
	return session.GoogleProvider.Get()
}
