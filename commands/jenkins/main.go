package main

import (
	"errors"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/application/property"
	"github.com/walterjwhite/go-application/libraries/external/jenkins"
)

var (
	jenkinsInstance *jenkins.Instance
	jenkinsJobFlag  = flag.String("j", "", "Jenkins Job Name")

	// move this to property/plugins/cli
	prefixFlag = flag.String("prefix", "", "property prefix, ie. if user specifies web/gmail.com/username with prefix of testing, resulting property would be testing/web/gmail.com/username")
)

func init() {
	application.Configure()

	jenkinsInstance = &jenkins.Instance{}

	property.Load(jenkinsInstance, *prefixFlag)

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
