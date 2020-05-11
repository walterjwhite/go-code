package jenkins

import (
	"github.com/walterjwhite/go-application/libraries/logging"

	"context"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-application/libraries/wait"
)

func (j *Job) Build(ctx context.Context) {
	j.get()

	var params map[string]string

	buildId, err := j.job.InvokeSimple(params)
	logging.Panic(err)

	log.Info().Msgf("%v - buildId: %v", j.job.GetName(), buildId)

	if j.BuildCheckInterval != nil && j.BuildTimeout != nil {
		wait.Wait(ctx, j.BuildCheckInterval, j.BuildTimeout, j.isDone)
	}
}
