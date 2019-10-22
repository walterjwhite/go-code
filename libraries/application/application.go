package application

import (
	"context"
	"flag"

	"github.com/walterjwhite/go-application/libraries/identifier"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/notification"
	"github.com/walterjwhite/go-application/libraries/shutdown"
)

var Context = shutdown.Default()

func Configure() context.Context {
	flag.Parse()

	logging.Configure()
	identifier.Log()

	return Context
}

func Wait() {
	// wait for CTRL+C (or context to expire)
	<-Context.Done()

	notification.OnCompletion()
}
