package jira

import (
	"github.com/walterjwhite/go-application/libraries/application/logging"
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func (i *Instance) Get(issueId string) *jiral.Issue {
	i.setupAuth()

	issue, _, err := i.client.Issue.Get(issueId, nil)
	logging.Panic(err)

	return issue
}
