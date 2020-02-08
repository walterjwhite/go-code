package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/property"
	"github.com/walterjwhite/go-application/libraries/walgreens"
)

var (
	session = &walgreens.Session{}
)

func init() {
	application.Configure()
	property.Load(session, "")

	log.Info().Msgf("session: %s", *session)
	property.Load(session.Credentials, "")
}

func main() {
	session.Authenticate(application.Context)

	application.Wait()
}
