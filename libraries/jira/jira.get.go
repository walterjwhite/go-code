package jira

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func (i *Instance) Get(issueId string) *jiral.Issue {
	issue, _, err := i.client.Issue.Get(issueId, nil)
	logging.Panic(err)

	return issue
}
