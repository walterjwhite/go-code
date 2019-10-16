package jenkins

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/wait"
)

func (j *JenkinsJob) Build(ctx context.Context) {
	var params map[string]string

	buildId, err := j.job.InvokeSimple(params)
	logging.Panic(err)

	log.Info().Msgf("%v - buildId: %v", j.job.GetName(), buildId)

	wait.Wait(ctx, j.jenkinsInstance.buildCheckInterval, j.jenkinsInstance.buildTimeout, j.isDone)
}

func (j *JenkinsJob) isDone() bool {
	running, err := j.job.IsRunning()
	logging.Panic(err)

	return !running
}
