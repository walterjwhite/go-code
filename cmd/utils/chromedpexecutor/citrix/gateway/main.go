package main

import (
	"errors"
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/cli"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/token/google"
	"time"
)

var (

	tickleInterval = flag.String("i", "3m", "Tickle Interval")
	session        = &citrix.Session{}
)

func init() {
	application.Configure(session)

	if len(*tickleInterval) > 0 {
		i, err := time.ParseDuration(*tickleInterval)
		logging.Panic(err)

		session.Tickle = &citrix.Tickle{TickleInterval: &i}
	}

	session.Validate()
	session.InitializeChromeDP(application.Context)
}

func main() {
	defer session.Session().Cancel()

	token := getToken()

	if !session.Run(*token) {
		logging.Panic(errors.New("unable to authenticate"))
	}

	session.Launch()
	application.Wait()
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
