package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/cli"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/google"
)

var (
	session = &citrix.Session{}
)

func init() {
	application.Configure(session)
	session.Validate()
	session.Init(application.Context)
}

func main() {
	if !session.Runnable() {
		log.Warn().Msg("Session will not start, past end time")
		return
	}

	defer session.Cancel()

	token := getToken()
	session.Run(*token)
}

func getToken() *string {
	token := cli.New().ReadToken(application.Context)
	if token != nil {
		log.Info().Msg("Using cmdline token")
		return token
	}

	googleProvider := &google.Provider{}
	application.Load(googleProvider)
	googleProvider.Init(application.Context)

	log.Info().Msgf("google: %v | %v | %v | %v", googleProvider.TokenTopicName, googleProvider.TokenSubscriptionName, googleProvider.StatusTopicName, googleProvider.StatusSubscriptionName)
	log.Info().Msgf("google Conf: %v | %v", googleProvider.Conf.CredentialsFile, googleProvider.Conf.ProjectId)

	return googleProvider.ReadToken(application.Context)
}
