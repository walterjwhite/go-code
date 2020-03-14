package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
)

func (j *Jira) Transition(taskName string, transition jiral.TransitionType) {
	jiral.Transition(j.jiraClient, j.IssueId, transition)
}
