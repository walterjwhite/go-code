package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/citrix/gateway"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/property"
	"time"
)

var (
	tokenFlag = flag.String("Token", "", "RSA Token")
	// TODO: randomize the interval, configure minimum interval and deviation ...
	tickleInterval = flag.String("TickleInterval", "1m", "Tickle Interval")
	session        = &gateway.Session{}
)

func init() {
	application.Configure()
	property.Load(session, "")

	log.Info().Msgf("session: %s", *session)
	property.Load(session.Credentials, "")
	log.Info().Msgf("session: %s", *session)

	i, err := time.ParseDuration(*tickleInterval)
	logging.Panic(err)

	session.Tickle = &gateway.Tickle{TickleInterval: &i}
}

func main() {
	session.Token = *tokenFlag
	session.Run(application.Context)

	application.Wait()
}
