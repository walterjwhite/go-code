package audit

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

func Audit(command *exec.Cmd, label string) (int, string, error) {
	logFile := path.GetFile(label, "log")

	var buffer bytes.Buffer

	err := ioutil.WriteFile(logFile.Name(), []byte(strings.Join(command.Args, " ")+"\n\n"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	runner.WithWriters(command, logFile, os.Stdout, &buffer)

	screenshot.Take(label, "0.before")
	err = command.Run()

	if err != nil {
		log.Printf("Error running command: %v\n", err)

		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), buffer.String(), err
		}
	}

	screenshot.Take(label, "1.after")
	return 0, buffer.String(), nil
}
