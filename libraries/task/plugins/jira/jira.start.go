package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
)

func (j *Jira) Start(taskName string) {
	jiral.Start(j.jiraClient, j.IssueId)
}
