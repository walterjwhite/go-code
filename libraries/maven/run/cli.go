package run

import (
	"context"
	"fmt"

	"log"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/io/writermatcher"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timestamp"
	"path/filepath"
)

func runApplication(ctx context.Context, index int, a Application) *exec.Cmd {
	log.Printf("Running Application: %v (%v)", a.Name, a.Command)
	log.Printf("Environment: %v", a.Environment)
	log.Printf("Arguments: %v", a.Arguments)
	log.Printf("Matcher: %v", a.LogMatcher)

	command := runner.Prepare(ctx, a.Command, a.Arguments...)
	notificationChannel := make(chan *string)

	logFile := getLogFile(application)
	runner.WithEnvironment(command, true, a.Environment...)

	a.configureLogWatcher(notificationChannel, logFile, command)

	logging.Panic(runner.Start(command))

	//go checkIfStarted(application, notificationChannel, notificationBuilder)

	return command
}

func (a *Application) configureLogWatcher(notificationChannel chan *string, logFile string, command *exec.Cmd) {
	if len(a.LogMatcher) > 0 {
		if "spring-boot" == a.LogMatcher {
			writer := writermatcher.NewSpringBootApplicationStartupMatcher(notificationChannel, &logFile)
			runner.WithWriter(command, writer)
		} else if "npm" == a.LogMatcher {
			writer := writermatcher.NewNPMStartupMatcher(notificationChannel, &logFile)
			runner.WithWriter(command, writer)
		} else {
			log.Printf("%v not matched, no log matcher configured.\n", a.LogMatcher)
		}
	}
}

func getLogFile(application string) *os.File {
	logFile := fmt.Sprintf("%s/.logs/%v", application, timestamp.Get())
	log.Printf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	logging.Panic(os.MkdirAll(filepath.Dir(logFile), os.ModePerm))

	outfile, err := os.Create(logFile)
	logging.Panic(err)

	return outfile
}
