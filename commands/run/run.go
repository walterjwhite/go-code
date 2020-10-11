package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	//"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/utils/run"
	//"github.com/walterjwhite/go-application/libraries/time/timeformatter/timestamp"

	"flag"
)

var (
	profileFlag = flag.String("p", "default", "profile")
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {
	// TODO: this is deprecated, integrate with task API
	//path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	i := run.New(flag.Args()...)
	i.Run(application.Context, *profileFlag)
}
