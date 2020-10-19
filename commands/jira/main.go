package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"github.com/walterjwhite/go-application/libraries/external/jira"
)

var (
	jiraInstance = &jira.Instance{}
)

func init() {
	application.ConfigureWithProperties(jiraInstance)

	validate()
}

func validate() {
	if len(flag.Args()) < 1 {
		logging.Panic(errors.New("Command is required (create, comment, transition, get)"))
	}
}

func main() {
	defer application.OnEnd()

	switch flag.Args()[0] {
	case "create":
		create()
	case "comment":
		comment()
	case "transition":
		transition()
	case "get":
		get()
	}
}
