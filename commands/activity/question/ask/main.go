package main

import (
	"github.com/walterjwhite/go-application/libraries/activity/plugins/question"
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
	if len(flag.Args()) > 0 {
		log.Info().Msgf("Raising Question: %v", flag.Args()[0])

		q := question.Ask(application.Context, flag.Args()[0])
		log.Info().Msgf("Raised Question: %v", q)
	} else {
		logging.Panic(fmt.Errorf("No question raised"))
	}
}
