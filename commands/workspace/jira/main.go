package main

import (
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task/plugins/jira"
)

var (
	actionFlag = flag.String("JiraAction", "comment", "comment on jira")
)

func init() {
	application.Configure()
}

// synchronous, waits for job to complete
func main() {
	defer application.OnEnd()

	t, name := plugins.InitializeWithName(application.Context)
	j := jira.Initialize(t, name)

	if *actionFlag == "comment" {
		j.Comment(flag.Args()[1])
	} else {
		j.Transition(flag.Args()[1])
	}
}
