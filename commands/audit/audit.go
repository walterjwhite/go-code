package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/audit"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"

	"errors"
	"flag"
)

var (
	scriptFlag = flag.String("RunScript", "", "Script file to run")
	labelFlag  = flag.String("Label", "", "Label to use")
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {
	if labelFlag == nil {
		logging.Panic(errors.New("Label is required"))
	}

	if scriptFlag != nil {
		audit.Run(application.Context, *scriptFlag, *labelFlag)
	} else {
		if len(flag.Args()) == 0 {
			logging.Panic(errors.New("Either specify a script file or command with arguments"))
		}

		command := flag.Args()[0]
		arguments := flag.Args()[1:]

		cmd := runner.Prepare(application.Context, command, arguments...)
		audit.Execute(cmd, *labelFlag)
	}
}
