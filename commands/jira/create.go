package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/external/jira"
)

var (
	createFlagSet = flag.NewFlagSet("create", flag.ExitOnError)

	projectKey  = createFlagSet.String("p", "", "Project Key")
	summary     = createFlagSet.String("s", "", "Summary")
	description = createFlagSet.String("d", "", "Description")
	issueType   = createFlagSet.String("t", "", "Issue Type")
)

func create() {
	logging.Panic(createFlagSet.Parse(flag.Args()[1:]))

	if len(*projectKey) == 0 {
		logging.Panic(errors.New("Project Key is required."))
	}

	if len(*summary) == 0 {
		logging.Panic(errors.New("Summary is required."))
	}

	if len(*description) == 0 {
		logging.Panic(errors.New("Description is required."))
	}

	if len(*issueType) == 0 {
		logging.Panic(errors.New("Issue Type is required."))
	}

	jira.Print(jiraInstance.Create(*projectKey, *summary, *description, *issueType))
}
