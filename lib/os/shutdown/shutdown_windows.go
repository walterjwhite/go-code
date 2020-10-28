package shutdown

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/utils/runner"
	"time"
)

func (r *ShutdownRequest) log() {
	arguments := r.getArguments()
	log.Info().Msgf("Arguments: %v", arguments)
}

func (r *ShutdownRequest) getArguments() []string {
	arguments := make([]string, 0)

	arguments = r.getShutdownAction(arguments)
	arguments = r.getTimeout(arguments)

	return arguments
}

func (r *ShutdownRequest) getShutdownAction(arguments []string) []string {
	if r.ShutdownAction == Reboot {
		arguments = append(arguments, "/r")
	} else if r.ShutdownAction == Poweroff {
		arguments = append(arguments, "/s")
	} else {
		logging.Panic(errors.New(fmt.Sprintf("Unknown option specified: %v\n", r.ShutdownAction)))
	}

	return arguments
}

func (r *ShutdownRequest) getTimeout(arguments []string) []string {
	arguments = append(arguments, "/t")
	arguments = append(arguments, fmt.Sprintf("%v", int64(r.Timeout/time.Second)))

	return arguments
}

func (r *ShutdownRequest) run() {
	arguments := r.getArguments()
	_, err := runner.Run(application.Context, "shutdown", arguments...)
	logging.Panic(err)
}
