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

// TODO: can this be done automatically without a direct mapping?
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
	logFile := filepath.Join(application, ".logs", timestamp.Get())
	log.Info().Msgf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	logging.Panic(os.MkdirAll(filepath.Dir(logFile), os.ModePerm))

	outfile, err := os.Create(logFile)
	logging.Panic(err)

	return outfile
}
