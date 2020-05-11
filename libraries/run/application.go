package run

import (
	"context"

	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/io/writermatcher"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timeformatter/timestamp"
	"path/filepath"
)

func (a *Application) Run(ctx context.Context, region string, index int) {
	a.command = runner.Prepare(ctx, a.Command, a.Arguments...)

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

	go monitorChannel(ctx, a.Name, notificationChannel)
}
