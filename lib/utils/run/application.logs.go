package run

import (
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/writermatcher"
	"github.com/walterjwhite/go-code/lib/time/timeformatter/timestamp"
	"github.com/walterjwhite/go-code/lib/utils/runner"
	"path/filepath"
)

// TODO: can this be done automatically without a direct mapping?
func (a *Application) configureLogWatcher(notificationChannel chan *string, writer io.Writer, command *exec.Cmd) {
	if len(a.LogMatcher) > 0 {
		// change this to be an external executable
		// pass stderr/stdout to this
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

func (a *Application) getLogFile() *os.File {
	logFile := filepath.Join(a.session.Path, a.Name, ".logs", timestamp.Get())
	log.Info().Msgf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	logging.Panic(os.MkdirAll(filepath.Dir(logFile), os.ModePerm))

	outfile, err := os.Create(logFile)
	logging.Panic(err)

	return outfile
}
