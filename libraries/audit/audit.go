package audit

import (
	"bytes"

	"os"
	"os/exec"
	"strings"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

//var auditPath = flag.String("AuditPath", "~/.audit", "AuditPath")

func init() {
	path.WithSessionDirectory("~/.audit")
}

// return code, output
func Execute(command *exec.Cmd, label string) (int, string) {
	logFile := path.GetFile(label, "log")
	defer logFile.Close()

	var buffer bytes.Buffer

	_, err := logFile.Write([]byte(strings.Join(command.Args, " ") + "\n\n"))
	logging.Panic(err)

	runner.WithWriters(command, logFile, os.Stdout, &buffer)

	screenshot.Take(label, "0.before")
	logging.Panic(command.Run())

	screenshot.Take(label, "1.after")

	return 0, buffer.String()
}
