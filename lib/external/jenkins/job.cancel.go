package jenkins

import (
	"github.com/walterjwhite/go/lib/application/logging"

	"context"
	"github.com/rs/zerolog/log"
)

func (j *Job) Cancel(ctx context.Context) {
	j.get()

	running, err := j.job.IsRunning()
	logging.Panic(err)

	if running {
		//success, err := j.GetBuild().Stop()
		//logging.Panic(err)
		build, err := j.job.GetLastBuild()
		logging.Panic(err)

		running = build.IsRunning()
		if running {
			running, err = build.Stop()
			logging.Panic(err)

			log.Info().Msgf("Stopped build %v / %v / %v", j.job.GetName(), build.GetBuildNumber(), running)
		} else {
			log.Info().Msgf("Build is not currently running")
		}
	} else {
		log.Info().Msgf("Job is not currently running")
	}
}
