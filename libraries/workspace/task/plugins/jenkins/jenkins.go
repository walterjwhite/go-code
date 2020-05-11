package jenkins

import (
	jenkinsl "github.com/walterjwhite/go-application/libraries/jenkins"
	"github.com/walterjwhite/go-application/libraries/property"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"github.com/walterjwhite/go-application/libraries/workspace/task/plugins"
)

func getJenkinsJob(t *task.Task, name string) *jenkinsl.Job {
	// TODO: implement

	// 1. store (jobName, buildTimeout, buildCheckInterval)
	// 2. reference jenkins instance (filename)
	// 3. decrypt secrets automatically

	var j *jenkinsl.Job
	plugins.Configure(t, name, j)

	// TODO: generalize this, support loading on nested fields
	if j.Instance == nil {
		// load default jenkins configuration
		property.Load(j.Instance, "")
	}

	return j
}
