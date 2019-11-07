package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/run"
	"github.com/walterjwhite/go-application/libraries/timestamp"

	"flag"
	"strings"
)

var (
	applications = flag.String("Applications", "default", "Comma-separated list of applications to run")
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {
	path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	run.Run(application.Context, getApplications(applications))
}

func getApplications(applicationsString *string) []string {
	return strings.Split(*applicationsString, ",")
}
