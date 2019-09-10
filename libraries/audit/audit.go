package audit

import (
	"bytes"
	"io/ioutil"

	"os"
	"os/exec"
	"strings"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

func init() {
	path.WithSessionDirectory("~.audit")
}

// return code, output
func Audit(command *exec.Cmd, label string) (int, string) {
	logFile := path.GetFile(label, "log")

	var buffer bytes.Buffer

	logging.Panic(ioutil.WriteFile(logFile.Name(), []byte(strings.Join(command.Args, " ")+"\n\n"), os.ModePerm))

	runner.WithWriters(command, logFile, os.Stdout, &buffer)

	screenshot.Take(label, "0.before")
	logging.Panic(command.Run())

	screenshot.Take(label, "1.after")
	return 0, buffer.String()
}
