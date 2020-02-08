package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/foreachfile"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/periodic"
	"github.com/walterjwhite/go-application/libraries/runner"

	"flag"
	"strings"
)

var (
	rootDirectoryFlag = flag.String("RootDirectory", ".", "Root Directory to scan files")
	intervalFlag      = flag.String("Interval", "1m", "Interval between execution")
	patternStringFlag = flag.String("Patterns", "", "Patterns")
)

func init() {
	application.Configure()
}

func main() {
	periodic.Periodic(application.Context, periodic.GetInterval(*intervalFlag), runIteration)

	application.Wait()
}

func runIteration() error {
	foreachfile.Execute(*rootDirectoryFlag, exec, getPatterns()...)

	return nil
}

func getPatterns() []string {
	if len(*patternStringFlag) > 0 {
		return strings.Split(*patternStringFlag, ",")
	}

	return []string{}
}

func exec(filePath string) {
	_, err := runner.Run(application.Context, filePath)
	logging.Panic(err)
}
