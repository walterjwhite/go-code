package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/activity/plugins/cli"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"strings"
	"time"
)

var (
	tagsFlag = flag.String("Tags", "", "Tags to apply")
	//envFlag        = flag.String("Env", "", "Env variables")
	dirFlag        = flag.String("Dir", "", "Dir to execute the command from")
	timeLimitFlag  = flag.String("TimeLimit", "1m", "Limit execution time")
	screenshotFlag = flag.Bool("Screenshots", false, "Take screenshots before and after")
)

func init() {
	application.Configure()
}

func main() {
	if len(flag.Args()) >= 1 {
		var env map[string]string

		tags := strings.Split(*tagsFlag, ",")

		timeLimit, err := time.ParseDuration(*timeLimitFlag)
		logging.Panic(err)

		cli.Execute(application.Context, flag.Args()[0], *dirFlag, env, timeLimit, *screenshotFlag, tags, flag.Args()[1:]...)
	} else {
		logging.Panic(fmt.Errorf("Unable to execute command, no command provided"))
	}
}
