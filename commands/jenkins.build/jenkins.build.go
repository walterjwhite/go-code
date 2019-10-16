package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/jenkins"
)

func main() {
	ctx := application.Configure()

	j := jenkins.New()
	job := j.GetCLIJob()
	job.Build(ctx)
}
