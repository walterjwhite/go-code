package main

import (
	"errors"

	"flag"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"os/exec"
	"strconv"
)

var (
	cmd = flag.String("cmd", "", "cmd to execute on receipt of pubsub message")
	arg = flag.String("arg", "", "arg to pass to cmd")
)

func init() {
	application.Configure()
}

func main() {
	if len(*cmd) == 0 {
		logging.Panic(errors.New("-cmd=<COMMAND>"))
	}

	if len(*arg) == 0 {
		logging.Panic(errors.New("-arg=<ARGUMENT>"))
	}

	log.Info().Msgf("running: %s -> %s", *cmd, *arg)
	ecmd := exec.Command(*cmd, *arg)
	output, err := ecmd.CombinedOutput()

	status := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
			log.Warn().Msgf("Error running: %s (%s) -> %v", *cmd, *arg, status)
		}
	}

	log.Info().Msgf("status: %v, output: %s", status, output)

	message := "status: " + strconv.Itoa(status) + "\noutput: " + string(output)
	log.Info().Msgf("message: %s", message)
}
