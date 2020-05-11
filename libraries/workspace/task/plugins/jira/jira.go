package jira

import (
	jiral "github.com/walterjwhite/go-application/libraries/jira"
	"github.com/walterjwhite/go-application/libraries/property"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"github.com/walterjwhite/go-application/libraries/workspace/task/plugins"
)

type Jira struct {
	UrlKey  string
	IssueId string

	Instance *jiral.Instance
}

func Initialize(t *task.Task, name string) *Jira {
	var j *Jira
	plugins.Configure(t, name, j)

	// load credentials specific to URLKey
	property.Load(j.Instance, j.UrlKey)

	return j
}
