package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os/user"
)

func init() {
	application.Configure()
}

func main() {
	currentUser, err := user.Current()
	logging.Panic(err)

	log.Info().Msgf("Current user: %v", currentUser.Username)
}
