package jenkins

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"time"
)

func (j *JenkinsJob) Build() {
	var params map[string]string

	buildId, err := j.job.InvokeSimple(params)
	logging.Panic(err)

	log.Info().Msgf("%v - buildId: %v", j.job.GetName(), buildId)
	logging.Panic(timeout.Limit(j.wait, j.jenkinsInstance.buildTimeout))
}

func (j *JenkinsJob) wait() {
	for {
		time.Sleep(j.jenkinsInstance.buildCheckInterval)
		running, err := j.job.IsRunning()
		logging.Panic(err)

		log.Info().Msgf("%v running: %v", j.job.GetName(), running)

		if !running {
			break
		}
	}
}
