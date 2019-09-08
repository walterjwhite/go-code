package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/audit"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timestamp"

	"os"
)

type NoArgumentsError struct {
}

func (e *NoArgumentsError) Error() string {
	return "No arguments passed"
}

func main() {
	ctx := application.Configure()

	path.WithSessionDirectory("~/.audit/" + timestamp.Get())

	if os.Args == nil || len(os.Args) == 1 {
		logging.Panic(&NoArgumentsError{})
	}

	command := os.Args[1]
	arguments := os.Args[2:]

	// log arguments
	cmd := runner.Prepare(ctx, command, arguments...)

	audit.Audit(cmd, "audit")
}
