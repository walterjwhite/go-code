package jenkins

import (
	//"github.com/bndr/gojenkins"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/bndr/gojenkins.v1"

	"context"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-application/libraries/wait"
)

type JenkinsJob struct {
	jenkinsInstance *JenkinsInstance
	job             *gojenkins.Job
}

func (i *JenkinsInstance) GetJob(jobName string) *JenkinsJob {
	i.setup()

	jobInstance, err := i.jenkins.GetJob(jobName)
	logging.Panic(err)

	return &JenkinsJob{job: jobInstance, jenkinsInstance: i}
}

func (j *JenkinsJob) Build(ctx context.Context) {
	var params map[string]string

	buildId, err := j.job.InvokeSimple(params)
	logging.Panic(err)

	log.Info().Msgf("%v - buildId: %v", j.job.GetName(), buildId)

	if j.jenkinsInstance.BuildCheckInterval != nil && j.jenkinsInstance.BuildTimeout != nil {
		wait.Wait(ctx, j.jenkinsInstance.BuildCheckInterval, j.jenkinsInstance.BuildTimeout, j.isDone)
	}
}

func (j *JenkinsJob) isDone() bool {
	running, err := j.job.IsRunning()
	logging.Panic(err)

	return !running
}
