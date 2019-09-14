package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/maven"
	"github.com/walterjwhite/go-application/libraries/maven/format"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/timestamp"

	"context"
	"flag"
	"fmt"
)

var debug = flag.Bool("Debug", false, "Whether maven should run with all the output or only WARN or higher")
var enabledLanguages = [len(format.Languages)]*bool{}
var hasLanguageEnabled = false

func init() {
	for i := 0; i < len(format.Languages); i++ {
		enabledLanguages[i] = flag.Bool(format.Languages[i].Name, false, fmt.Sprintf("Format only %v code\n", format.Languages[i].Name))
	}
}

// TODO: integrate win10 / dbus notifications
func main() {
	ctx := application.Configure()

	path.WithSessionDirectory("~/.audit/maven/format/" + timestamp.Get())

	process(ctx)
}

func process(ctx context.Context) {
	for i := 0; i < len(format.Languages); i++ {
		if *enabledLanguages[i] {
			hasLanguageEnabled = true
		}
	}

	for i := 0; i < len(format.Languages); i++ {
		if !hasLanguageEnabled || *enabledLanguages[i] {
			command := format.Languages[i].Command
			arguments := format.Languages[i].Arguments

			if !*debug {
				_, arguments = maven.GetCommandLine(arguments, debug)
			}

			_, err := runner.Run(ctx, command, arguments...)
			logging.Panic(err)
		}
	}
}
