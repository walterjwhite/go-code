package jenkins

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/until"
)

func (j *Job) Until(ctx context.Context) {
	if j.Instance.BuildCheckInterval != nil && j.Instance.BuildTimeout != nil {
		until.New(ctx, j.Instance.BuildCheckInterval, j.Instance.BuildTimeout, j.isDone)
	} else {
		log.Warn().Msg("Build check interval or build timeout is nil")
	}
}
