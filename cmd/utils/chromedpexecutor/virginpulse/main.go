package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application"

	"github.com/walterjwhite/go/lib/application/property"
	"github.com/walterjwhite/go/lib/virginpulse"
)

var (
	session = &virginpulse.Session{}
	//Credentials: &virginpulse.Credentials{}}
)

func init() {
	application.Configure()

	property.Load(session, "")

	log.Info().Msgf("session: %v", *session)
	property.Load(session.Credentials, "")
	log.Info().Msgf("session: %v", *session)
}

func main() {
	defer application.OnEnd()

	session.Authenticate(application.Context)
	session.RunScript()
	session.Logout()
}
