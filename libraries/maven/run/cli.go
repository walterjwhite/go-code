package run

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"strings"

	"github.com/walterjwhite/go-application/libraries/io/writermatcher"
	"github.com/walterjwhite/go-application/libraries/notify"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timestamp"
	"path/filepath"
)

func runApplication(ctx context.Context, index int, profile string, configuration Configuration, application string, debug bool, notificationBuilder func(notification notify.Notification) notify.Notifier) *exec.Cmd {
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

	var writer io.Writer = logFilter.NewSpringBootWriterFilter(notificationChannel, logFile)
	runner.WithWriter(command, &writer)

	err2 := runner.Start(command)
	if err2 != nil {
		log.Fatalf("Error running %v: %v\n", command, err2)
	}

	go checkIfStarted(application, notificationChannel, notificationBuilder)

	return command
}

func getDebugPort(index int) int {
	return DebugPortStart + index
}

func getJarFile(application string) *string {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/target", application))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Index(f.Name(), "jar") >= 0 && strings.Index(f.Name(), ".original") == -1 {
			jarFile := fmt.Sprintf("%s/target/%s", application, f.Name())
			return &jarFile
		}
	}

	log.Fatalf("No application artifact found for %s\n", application)
	return nil
}

func getLogFile(application string) *os.File {
	logFile := fmt.Sprintf("%s/target/logs/%v", application, timestamp.Get())
	log.Printf("writing logs to: %s", logFile)

	return makeLogFile(logFile)
}

func makeLogFile(logFile string) *os.File {
	os.MkdirAll(filepath.Dir(logFile), os.ModePerm)

	outfile, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}

	return outfile
}
