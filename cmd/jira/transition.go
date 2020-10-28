package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/external/jira"
)

var (
	transitionFlagSet = flag.NewFlagSet("transition", flag.ExitOnError)

	transitionIssueId = transitionFlagSet.String("i", "", "Issue Id")
	transitionAction  = transitionFlagSet.String("a", "", "Transition Action")
)

func transition() {
	logging.Panic(transitionFlagSet.Parse(flag.Args()[1:]))

	if len(*transitionIssueId) == 0 {
		logging.Panic(errors.New("Issue Id is required."))
	}

	if len(*transitionAction) == 0 {
		logging.Panic(errors.New("Transition Action is required."))
	}

	jiraInstance.Transition(*transitionIssueId, *transitionAction)
	jira.Print(jiraInstance.Get(*transitionIssueId))
}
