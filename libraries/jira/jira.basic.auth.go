package jira

import (
	"github.com/walterjwhite/go-application/libraries/property"
	jiral "gopkg.in/andygrunwald/go-jira.v1"
)

func NewBasicAuthClient() *Instance {
	i := &Instance{Credentials: &Credentials{}}

	property.Load(i.Credentials, "")
	transport := jiral.BasicAuthTransport{Username: i.Credentials.Username, Password: i.Credentials.Password}
	jiraClient, _ := jiral.NewClient(transport.Client(), i.Uri)

	i.client = jiraClient
	return i
}
