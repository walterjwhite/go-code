package main

import (
	"github.com/walterjwhite/go-application/libraries/activity"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"

	"errors"
	"flag"
)

var (
	scriptFlag = flag.String("RunScript", "", "Script file to run")
	labelFlag  = flag.String("Label", "", "Label to use")
)

func init() {
	application.Configure()
}

// TODO: integrate win10 / dbus notifications
func main() {

}
