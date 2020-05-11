package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/audit"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) == 0 {
		logging.Panic(errors.New("No arguments passed"))
	}

	command := flag.Args()[0]
	arguments := flag.Args()[1:]

	// log arguments
	cmd := runner.Prepare(application.Context, command, arguments...)

	audit.Execute(cmd, "audit")
}
