package main

import (
	//"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/os/user/idle"
)

func init() {
	application.Configure()
}

func main() {
	log.Info().Msgf("User idle time: %v", idle.IdleTime(application.Context))
}
