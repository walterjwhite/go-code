package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	//"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/run"
	//"github.com/walterjwhite/go-application/libraries/timeformatter/timestamp"

	"flag"
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {
	// TODO: this is deprecated, integrate with task API
	//path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	run.Run(application.Context, flag.Args())
}
