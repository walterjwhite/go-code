package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/bnymellon401k"
)

var (
	session = &bnymellon401k.Session{}
)

func init() {
	application.ConfigureWithProperties(session)

	log.Info().Msgf("Username: %s", session.Credentials.Username)
	log.Info().Msgf("Password: %s", session.Credentials.Password)
}

func main() {
	defer application.OnEnd()

	session.Login(context.Background())
}
