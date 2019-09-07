package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Audit(command *exec.Cmd, label string) (int, string, error) {
	logFile := path.GetFile(label, "log")

	var buffer bytes.Buffer

	ioutil.WriteFile(logFile.Name(), []byte(strings.Join(command.Args, " ")+"\n\n"), os.ModePerm)

	runner.WithWriters(command, logFile, os.Stdout, &buffer)

	screenshot.Screenshot(label, "0.before")
	err := command.Run()

	if err != nil {
		log.Printf("Error running command: %v\n", err)

		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), buffer.String(), err
		}
	}

	screenshot.Screenshot(label, "1.after")
	return 0, buffer.String(), nil
}
