package jenkins

import (
	"gopkg.in/bndr/gojenkins.v1"
	"time"
)

type Job struct {
	Name               string
	BuildTimeout       *time.Duration
	BuildCheckInterval *time.Duration

	Instance *Instance
	job      *gojenkins.Job
}
