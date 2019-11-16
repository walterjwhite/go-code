package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/jenkins"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/property"
)

var (
	jenkinsInstance *jenkins.JenkinsInstance
	jenkinsJobFlag  = flag.String("JenkinsJobName", "", "Jenkins Job Name")
)

func init() {
	application.Configure()

	jenkinsInstance = &jenkins.JenkinsInstance{}

	property.Load(jenkinsInstance, "")

	validate()
}

func main() {
	job := jenkinsInstance.GetJob(*jenkinsJobFlag)
	job.Build(application.Context)
}

func validate() {
	if len(*jenkinsJobFlag) == 0 {
		logging.Panic(errors.New("Jenkins Job Name is required."))
	}
}
