package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
)

type Jira struct {
	Url string
	IssueId string
	
	credentials *jiral.Credentials
	jiraClient *jirac.Client
}

func initialize() *Jira {
	// load credentials specific to URL
	property.Load(j.credentials, j.Url)
	
	auth := jira.BasicAuthTransport{Username: j.credentials.Username, Password: j.credentials.Password}
	jiraClient, _ := jira.NewClient(tp.Client(), <url>)
}
