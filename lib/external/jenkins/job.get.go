package jenkins

import (
	"github.com/walterjwhite/go/lib/application/logging"
)

func (i *Instance) GetJob(jobName string) *Job {
	j := &Job{Name: jobName, Instance: i}
	j.get()

	return j
}

func (j *Job) get() {
	if j.Instance.jenkins == nil {
		j.Instance.setup()
	}

	jobInstance, err := j.Instance.jenkins.GetJob(j.Name)
	logging.Panic(err)

	j.job = jobInstance
}
