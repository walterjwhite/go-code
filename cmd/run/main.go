package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/utils/run"
)

var (
	sessionPathFlag = flag.String("session-path", "", "session path (where app confs are located)")
)

func init() {
	application.Configure()
}

// TODO: integrate notifications
func main() {
	if len(*sessionPathFlag) == 0 {
		logging.Panic(errors.New("session path is required"))
	}

	if len(flag.Args()) == 0 {
		logging.Panic(errors.New("At least 1 application to run is required"))
	}

	s := run.New(*sessionPathFlag, flag.Args()...)
	s.Run(application.Context)
}
