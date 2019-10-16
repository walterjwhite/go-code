package jenkins

import (
	"flag"
	//"github.com/bndr/gojenkins"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/bndr/gojenkins.v1"
)

var jenkinsJobName = flag.String("JenkinsJobName", "", "JenkinsJobName")

type JenkinsJob struct {
	jenkinsInstance *JenkinsInstance
	job             *gojenkins.Job
}

type NoJobNameSpecifiedError struct{}

func (e *NoJobNameSpecifiedError) Error() string {
	return "No Job Name was specified for the option: JenkinsJobName"
}

func (j *JenkinsInstance) GetCLIJob() *JenkinsJob {
	if len(*jenkinsJobName) == 0 {
		logging.Panic(&NoJobNameSpecifiedError{})
	}

	return j.GetJob(*jenkinsJobName)
}

func (j *JenkinsInstance) GetJob(jobName string) *JenkinsJob {
	jobInstance, err := j.jenkins.GetJob(jobName)
	logging.Panic(err)

	return &JenkinsJob{job: jobInstance, jenkinsInstance: j}
}
