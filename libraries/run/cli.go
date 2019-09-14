package run

import (
	"context"
	"fmt"

	"io"
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
	command := runner.Prepare(ctx, a.Command, a.Arguments...)

	command.Dir = a.Name

	notificationChannel := make(chan *string)

	logFile := getLogFile(a.Name)
	runner.WithEnvironment(command, true, a.Environment...)

	a.configureLogWatcher(notificationChannel, logFile, command)

	log.Printf("Running Application: %v (%v)", a.Name, a.Command)
	log.Printf("Environment: %v", a.Environment)
	log.Printf("Arguments: %v", a.Arguments)
	log.Printf("Matcher: %v", a.LogMatcher)

	logging.Panic(runner.Start(command))

	go monitorChannel(a.Name, notificationChannel)

	return command
}

func (a *Application) configureLogWatcher(notificationChannel chan *string, writer io.Writer, command *exec.Cmd) {
	if len(a.LogMatcher) > 0 {
		if "spring-boot" == a.LogMatcher {
			writer = writermatcher.NewSpringBootApplicationStartupMatcher(notificationChannel, writer)
		} else if "npm" == a.LogMatcher {
			writer = writermatcher.NewNPMStartupMatcher(notificationChannel, writer)
		} else {
			log.Printf("%v not matched, no log matcher configured.\n", a.LogMatcher)
		}
	}

	runner.WithWriter(command, writer)
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
