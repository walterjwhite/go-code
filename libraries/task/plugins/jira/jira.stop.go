package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
)

func (j *Jira) Stop(taskName string) {
	jiral.Stop(j.jiraClient, j.IssueId)
}
