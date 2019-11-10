package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	// TODO: support passing in a specific hash to edit ...
	//commitHashFlag = flag.String("CommitHash", "", "Commit Hash")
	dateString string
)

func init() {
	application.Configure()

	dateString = flag.Args()[0]
	if len(dateString) == 0 {
		logging.Panic(errors.New("Date is required: ie. Wed Dec 19 09:21:37 2018 -0500"))
	}
}

// TODO: integrate win10 / dbus notifications
func main() {
	//path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	git.AmendDate(application.Context, dateString)
}
