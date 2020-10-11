package main

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins/jenkins"
)

var (
	cancelFlag = flag.Bool("Cancel", false, "Cancel the job (if it is running)")
)

func init() {
	application.Configure()
}

// synchronous, waits for job to complete
func main() {
	defer application.OnEnd()

	t, name := plugins.InitializeWithName(application.Context)

	if *cancelFlag {
		jenkins.Cancel(application.Context, t, name)
	} else {
		jenkins.Build(application.Context, t, name)
	}
}
