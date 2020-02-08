package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/activity"
	"github.com/walterjwhite/go-application/libraries/activity/plugins/cli"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"strings"
	"time"
)

var (
	tagsFlag = flag.String("Tags", "", "Tags to apply")
	//envFlag        = flag.String("Env", "", "Env variables")
	dirFlag         = flag.String("Dir", "", "Dir to execute the command from")
	timeLimitFlag   = flag.String("TimeLimit", "1m", "Limit execution time")
	screenshotFlag  = flag.Bool("Screenshots", false, "Take screenshots before and after")
	baseProjectPath = flag.String("BaseProjectPath", "ssh://git@localhost:/projects/active/", "Base Project Path")
	scriptFileFlag  = flag.String("ScriptFile", "", "Script File to run")
)

func init() {
	application.Configure()
}

func main() {
	if len(*scriptFileFlag) > 0 {
		var scriptFile *cli.ScriptFile
		yamlhelper.Read(*scriptFileFlag, scriptFile)

		cli.Run(application.Context, scriptFile)
	} else if len(flag.Args()) >= 1 {
		var env map[string]string

		tags := strings.Split(*tagsFlag, ",")

		timeLimit, err := time.ParseDuration(*timeLimitFlag)
		logging.Panic(err)

		activity.BaseProjectLocation = *baseProjectPath

		command := cli.New(flag.Args()[0], *dirFlag, env, &timeLimit, *screenshotFlag, tags, flag.Args()[1:]...)
		cli.Execute(application.Context, command)
	} else {
		logging.Panic(fmt.Errorf("Unable to execute command, no command provided"))
	}
}
