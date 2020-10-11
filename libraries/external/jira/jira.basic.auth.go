package jira

import (
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func (i *Instance) setupAuth() {
	if i.client != nil {
		return
	}

	transport := jiral.BasicAuthTransport{Username: i.Credentials.Username, Password: i.Credentials.Password}
	jiraClient, _ := jiral.NewClient(transport.Client(), i.Uri)

	i.client = jiraClient
}
