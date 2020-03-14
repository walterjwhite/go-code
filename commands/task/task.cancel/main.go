package main

import (
	"github.com/walterjwhite/go-application/libraries/task"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"fmt"
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) == 1 {
		log.Info().Msgf("Canceling task: %v", flag.Args()[0])

		task.Cancel(application.Context, flag.Args()[0])
	} else {
		logging.Panic(fmt.Errorf("Path for task is required ..."))
	}
}
