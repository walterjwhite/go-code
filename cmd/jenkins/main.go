package main

import (
	"errors"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/external/jenkins"
)

var (
	jenkinsInstance *jenkins.Instance
	jenkinsJobFlag  = flag.String("j", "", "Jenkins Job Name")
)

func init() {
	jenkinsInstance = &jenkins.Instance{}

	application.ConfigureWithProperties(jenkinsInstance)

	validate()

	log.Info().Msgf("Looking for job: %v", *jenkinsJobFlag)
}

func validate() {
	if len(flag.Args()) < 1 {
		logging.Panic(errors.New("Command is required (build, cancel, wait)"))
	}

	if len(*jenkinsJobFlag) == 0 {
		logging.Panic(errors.New("Jenkins Job Name is required."))
	}
}

func main() {
	defer application.OnEnd()

	job := jenkinsInstance.GetJob(*jenkinsJobFlag)

	switch flag.Args()[0] {
	case "build":
		job.Build(application.Context)
	case "cancel":
		job.Cancel(application.Context)
	case "wait":
		job.Wait(application.Context)
	}
}
