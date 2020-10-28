package jenkins

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/time/wait"
)

func (j *Job) Wait(ctx context.Context) {
	if j.Instance.BuildCheckInterval != nil && j.Instance.BuildTimeout != nil {
		wait.Wait(ctx, j.Instance.BuildCheckInterval, j.Instance.BuildTimeout, j.isDone)
	} else {
		log.Warn().Msg("Build check interval or build timeout is nil")
	}
}
