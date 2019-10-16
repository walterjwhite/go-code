package jenkins

import (
	//"github.com/bndr/gojenkins"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/bndr/gojenkins.v1"
)

type JenkinsJob struct {
	jenkinsInstance *JenkinsInstance
	job             *gojenkins.Job
}

func (j *JenkinsInstance) GetJob(jobName string) *JenkinsJob {
	jobInstance, err := j.jenkins.GetJob(jobName)
	logging.Panic(err)

	return &JenkinsJob{job: jobInstance, jenkinsInstance: j}
}
