package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/jenkins"
)

func init() {
	application.Configure()
}

func main() {
	j := jenkins.New()
	job := j.GetCLIJob()
	job.Build(application.Context)
}
