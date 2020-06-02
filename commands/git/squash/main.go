package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func init() {
	application.Configure()

	if len(flag.Args()) != 1 {
		logging.Panic(errors.New("Expecting exactly 1 argument, the number of commits to squash"))
	}
}

func main() {
	defer application.OnEnd()

	status, err := runner.Run(application.Context, "git", "reset", "--soft", fmt.Sprintf("HEAD~%s", flag.Args()[0]))
	logging.Panic(err)

	if status > 0 {
		logging.Panic(fmt.Errorf("Expecting exit status to be 0 %d", status))
	}

	runner.Run(application.Context, "git", "commit")
}
