package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/maven/build"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/timestamp"

	"flag"
)

var debug = flag.Bool("Debug", false, "Whether maven should run with all the output or only WARN or higher")

// TODO: integrate win10 / dbus notifications
func main() {
	ctx := application.Configure()

	path.WithSessionDirectory("~/.audit/maven/build/" + timestamp.Get())

	build.Build(ctx, debug)
}
