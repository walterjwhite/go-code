package application

import (
	"context"
	"flag"

	"github.com/walterjwhite/go-application/libraries/identifier"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/notification"
	"github.com/walterjwhite/go-application/libraries/shutdown"
)

var logFile = flag.String("Log", "", "The log file to write to")
var Context = shutdown.Default()

func Configure() context.Context {
	identifier.Log()

	flag.Parse()

	logging.Set(*logFile)

	return Context
}

func OnCompletion() {
	notification.OnCompletion()
}

func Wait() {
	// wait for CTRL+C (or context to expire)
	<-Context.Done()
}
