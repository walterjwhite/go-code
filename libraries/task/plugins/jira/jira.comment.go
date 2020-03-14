package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
)

func (j *Jira) Comment(taskName, comment string) {
	jiral.Comment(j.jiraClient, j.IssueId, comment)
}
