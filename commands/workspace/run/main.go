package main

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins/run"
)

func init() {
	application.Configure()
}

// synchronous, waits for job to complete
func main() {
	defer application.OnEnd()

	t := plugins.Initialize(application.Context)

	if len(flag.Args()) >= 1 {
		run.Run(application.Context, t, flag.Args()...)
	} else {
		run.Run(application.Context, t, "")
	}
}
