package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/jira"
)

var (
	getFlagSet = flag.NewFlagSet("get", flag.ExitOnError)

	getIssueId = getFlagSet.String("i", "", "Issue Id")
)

func get() {
	logging.Panic(getFlagSet.Parse(flag.Args()[1:]))

	if len(*getIssueId) == 0 {
		logging.Panic(errors.New("Issue Id is required."))
	}

	jira.Print(jiraInstance.Get(*getIssueId))
}
