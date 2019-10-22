package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/screenshot"
)

var label = flag.String("Label", "", "Screenshot label")
var detail = flag.String("Detail", "", "Screenshot detail")

func main() {
	application.Configure()

	if len(*label) == 0 {
		logging.Panic(errors.New("Please specify a label for the screenshot"))
	}

	if len(*detail) == 0 {
		logging.Panic(errors.New("Please specify a detailed message for the screenshot"))
	}

	path.WithSessionDirectory("~/.audit")
	screenshot.Take(*label, *detail)
}
