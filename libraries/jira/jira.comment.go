package jira

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func (i *Instance) Comment(issueId, comment string) *jiral.Issue {
	issue := i.Get(issueId)

	c := &jiral.Comment{Body: comment}
	_, _, err := i.client.Issue.AddComment(issue.ID, c)
	logging.Panic(err)

	return issue
}
