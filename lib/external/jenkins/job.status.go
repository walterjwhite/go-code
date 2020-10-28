package jenkins

import (
	"github.com/walterjwhite/go/lib/application/logging"
)

func (j *Job) isDone() bool {
	running, err := j.job.IsRunning()
	logging.Panic(err)

	return !running
}
