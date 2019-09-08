package run

import (
	"context"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"strings"

	"github.com/walterjwhite/go-application/libraries/io/writermatcher"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timestamp"
	"path/filepath"
)

func runApplication(ctx context.Context, index int, profile string, configuration Configuration, application string, debug bool /*, notificationBuilder func(notification notify.Notification) notify.Notifier*/) *exec.Cmd {
	log.Printf("Running Application: %v (%v)", application, profile)
	log.Printf("Environment: %v", configuration.Environment)
	log.Printf("JVMArguments: %v", configuration.Jvm)

	arguments := make([]string, 0)
	if debug {
		arguments = append(arguments, fmt.Sprintf(DebugArguments, getDebugPort(index)))
	}

	for _, jvmArgument := range configuration.Jvm {
		arguments = append(arguments, fmt.Sprintf("-D%v", jvmArgument))
	}

	arguments = append(arguments, "-jar")
	arguments = append(arguments, *getJarFile(application))

	command := runner.Prepare(ctx, "java", arguments...)
	notificationChannel := make(chan *string)

	logFile := getLogFile(application)
	runner.WithEnvironment(command, true, configuration.Environment...)

	// TODO: configure this @ runtime, perhaps we're not using SpringBoot @ all
	writer := writermatcher.NewSpringBootApplicationStartupMatcher(notificationChannel, logFile)
	runner.WithWriter(command, writer)

	logging.Panic(runner.Start(command))

	//go checkIfStarted(application, notificationChannel, notificationBuilder)

	return command
}

func getDebugPort(index int) int {
	return DebugPortStart + index
}

type NoApplicationArtifactFoundError struct {
	ApplicationName string
}

func (e *NoApplicationArtifactFoundError) Error() string {
	return fmt.Sprintf("No application artifact found for %s\n", e.ApplicationName)
}

func getJarFile(application string) *string {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/target", application))
	logging.Panic(err)

	for _, f := range files {
		if strings.Contains(f.Name(), "jar") && !strings.Contains(f.Name(), ".original") {
			jarFile := fmt.Sprintf("%s/target/%s", application, f.Name())
			return &jarFile
		}
	}

	logging.Panic(&NoApplicationArtifactFoundError{ApplicationName: application})
	return nil
}

func getLogFile(application string) *os.File {
	logFile := fmt.Sprintf("%s/target/logs/%v", application, timestamp.Get())
	log.Printf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	logging.Panic(os.MkdirAll(filepath.Dir(logFile), os.ModePerm))

	outfile, err := os.Create(logFile)
	logging.Panic(err)

	return outfile
}
