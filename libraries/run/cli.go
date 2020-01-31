package run

import (
	"context"
	"fmt"

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

func runApplication(ctx context.Context, index int, a Application) *exec.Cmd {
	command := runner.Prepare(ctx, a.Command, a.Arguments...)

	command.Dir = a.Name

	notificationChannel := make(chan *string)

	logFile := getLogFile(a.Name)
	runner.WithEnvironment(command, true, a.Environment...)

	a.configureLogWatcher(notificationChannel, logFile, command)

	log.Info().Msgf("Running Application: %v (%v)", a.Name, a.Command)
	log.Info().Msgf("Environment: %v", a.Environment)
	log.Info().Msgf("Arguments: %v", a.Arguments)
	log.Info().Msgf("Matcher: %v", a.LogMatcher)

	logging.Panic(command.Start())

	go monitorChannel(ctx, a.Name, notificationChannel)

	return command
}

func (a *Application) configureLogWatcher(notificationChannel chan *string, writer io.Writer, command *exec.Cmd) {
	if len(a.LogMatcher) > 0 {
		if "spring-boot" == a.LogMatcher {
			writer = writermatcher.NewSpringBootApplicationStartupMatcher(notificationChannel, writer)
		} else if "npm" == a.LogMatcher {
			writer = writermatcher.NewNPMStartupMatcher(notificationChannel, writer)
		} else {
			log.Info().Msgf("%v not matched, no log matcher configured.\n", a.LogMatcher)
		}
	}

	runner.WithWriter(command, writer)
}

func getLogFile(application string) *os.File {
	logFile := fmt.Sprintf("%s/.logs/%v", application, timestamp.Get())
	log.Info().Msgf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	logging.Panic(os.MkdirAll(filepath.Dir(logFile), os.ModePerm))

	outfile, err := os.Create(logFile)
	logging.Panic(err)

	return outfile
}
