package main

import (
	"flag"

	"errors"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/external/jira"
)

var (
	commentFlagSet = flag.NewFlagSet("comment", flag.ExitOnError)

	commentIssueId = commentFlagSet.String("i", "", "Issue ID")
	commentText    = commentFlagSet.String("c", "", "Comment")
)

func comment() {
	logging.Panic(commentFlagSet.Parse(flag.Args()[1:]))

	if len(*commentIssueId) == 0 {
		logging.Panic(errors.New("Issue ID is required."))
	}

	if len(*commentText) == 0 {
		logging.Panic(errors.New("Comment is required."))
	}

	jira.Print(jiraInstance.Comment(*commentIssueId, *commentText))
}
