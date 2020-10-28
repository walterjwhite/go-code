package jenkins

import (
	"gopkg.in/bndr/gojenkins.v1"
)

type Job struct {
	Name string

	Instance *Instance
	job      *gojenkins.Job
}
