package main

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/screenshot"
	"log"
)

var label = flag.String("Label", "", "Screenshot label")
var detail = flag.String("Detail", "", "Screenshot detail")

func main() {
	application.Configure()
	defer application.OnCompletion()

	if len(*label) == 0 {
		log.Fatal("Please specify a label for the screenshot")
	}

	if len(*detail) == 0 {
		log.Fatal("Please specify a detailed message for the screenshot")
	}

	path.WithSessionDirectory("~/.audit")
	screenshot.Take(*label, *detail)
}
