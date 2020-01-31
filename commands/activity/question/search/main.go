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
		log.Info().Msgf("Searching for question: %v", flag.Args()[0])

		r := question.Search(application.Context, flag.Args()[0])
		log.Info().Msgf("Found %v Questions:", len(r))

		for i, q := range r {
			log.Info().Msgf("Question(%v): %v: %v", i, *q)
		}
	} else {
		logging.Panic(fmt.Errorf("No search text entered"))
	}
}
