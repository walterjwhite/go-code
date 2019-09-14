package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/maven/run"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/timestamp"

	"flag"
)

var applications = flag.String("Applications", "default", "Comma-separated list of applications to run")

// TODO: integrate win10 / dbus notifications
func main() {
	ctx := application.Configure()

	path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	run.Run(ctx, getApplications(applications))
}

func getApplications(applicationsString *string) {
	return strings.Split(*applicationsString, ",")
}
