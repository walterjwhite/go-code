package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/cli"
)

var (
	session       = &citrix.Session{}
	firstRun bool = true
)

func init() {
	application.Configure(session)
	session.Validate()

	log.Debug().Msg("initializing google pubsub")

	application.Load(session.GoogleProvider)
	session.GoogleProvider.Init(application.Context)

	log.Debug().Msgf("conf: %v", session.GoogleProvider)
}

func main() {
	for {
		token := getToken()

		session.Init(application.Context)
		session.Run(*token)

		firstRun = false
	}
}

func getToken() *string {
	if firstRun {
		token := cli.New().ReadToken()
		if token != nil {
			log.Info().Msg("using cmdline token")
			return token
		}
	}

	log.Info().Msg("using google pubsub")
	return session.GoogleProvider.ReadToken()
}
