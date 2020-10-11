package jenkins

import (
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"context"
	"github.com/rs/zerolog/log"
)

func (j *Job) Build(ctx context.Context) {
	// j.get()

	var params map[string]string

	buildId, err := j.job.InvokeSimple(params)
	logging.Panic(err)

	log.Info().Msgf("%v - buildId: %v", j.job.GetName(), buildId)

	j.Wait(ctx)
}
