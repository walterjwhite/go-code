package run

import (
	"context"
	"os/exec"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/utils/runner"
)

func (a *Application) Run(ctx context.Context, region string, index int) {
	a.command = exec.CommandContext(ctx, a.Command, a.Arguments...)

	a.command.Dir = a.Name

	notificationChannel := make(chan *string)

	logFile := getLogFile(a.Name)
	runner.WithEnvironment(a.command, true, a.Environment...)

	a.configureLogWatcher(notificationChannel, logFile, a.command)

	log.Info().Msgf("Running Application: %v (%v)", a.Name, a.Command)
	log.Info().Msgf("Environment: %v", a.Environment)
	log.Info().Msgf("Arguments: %v", a.Arguments)
	log.Info().Msgf("Matcher: %v", a.LogMatcher)

	logging.Panic(a.command.Start())

	go a.monitorChannel(ctx, notificationChannel)
	go a.checkPort(ctx)
}
