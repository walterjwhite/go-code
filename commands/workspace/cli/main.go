package main

import (
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins/cli"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"flag"
	"fmt"
)

var (
	captureScreenshotsFlag = flag.Bool("CaptureScreenshots", true, "Take screenshots before and after execution (defaults to yes)")
)

func init() {
	application.Configure()
}

// TODO:
// 1. use task.plugin helper to locate task directory
// 2. with task directory, determine job details from <name>.jenkins configuration file

// synchronous, waits for job to complete
func main() {
	defer application.OnEnd()

	if len(flag.Args()) > 1 {
		t := plugins.Initialize(application.Context)

		cli.Execute(application.Context, t, *captureScreenshotsFlag, flag.Args()[0], flag.Args()[1:]...)
	} else {
		logging.Panic(fmt.Errorf("Cmd is required, arguments are optional ..."))
	}
}
