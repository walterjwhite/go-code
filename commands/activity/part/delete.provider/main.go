package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/part"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 2 {
		log.Info().Msgf("Deleting Part Provider: %v", flag.Args()[0])

		part.DeleteProvider(application.Context, flag.Args()[0], flag.Args()[1])
		log.Info().Msgf("Deleted Part Provider: %v", flag.Args()[1])
	} else {
		logging.Panic(fmt.Errorf("No provider deleted"))
	}
}
