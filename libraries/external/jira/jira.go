package jira

import (
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

type Credentials struct {
	Username string
	Password string
}

func (i *Instance) SecretFields() []string {
	return []string{"Credentials.Username", "Credentials.Password"}
}

type Instance struct {
	Credentials             *Credentials
	Url                     string
	TransitionActionMapping map[string]int

	client *jiral.Client
}

/*
type IssueType string
const ()

i := get(jiraClient)
comment(jiraClient, i, "updated - WW.2020/03/12")
create(jiraClient)

transition(jiraClient, ticketId, StartProgress)
transition(jiraClient, ticketId, StopProgress)
transition(jiraClient, ticketId, DevComplete)
transition(jiraClient, ticketId, ReadyForQA)
transition(jiraClient, ticketId, Close)
transition(jiraClient, ticketId, Reopen)
*/
