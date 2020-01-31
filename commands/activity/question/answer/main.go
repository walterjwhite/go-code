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
	if len(flag.Args()) == 2 {
		q := question.DoAnswer(application.Context, flag.Args()[0], flag.Args()[1])
		log.Info().Msgf("Answered Question: %v", q)
	} else {
		logging.Panic(fmt.Errorf("Question id and Answer text are both required."))
	}
}
