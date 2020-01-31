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
	if len(flag.Args()) > 1 {
		log.Info().Msgf("Merging for questions: %v", flag.Args())

		q := question.DoMerge(application.Context, flag.Args()...)
		log.Info().Msgf("Merged into Question:", *q)
	} else {
		logging.Panic(fmt.Errorf("2 or more arguments are required to merge"))
	}
}
