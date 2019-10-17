package jenkins

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"gopkg.in/bndr/gojenkins.v1"
	"time"
)

type JenkinsInstance struct {
	url      string
	username string
	password string

	buildTimeout       time.Duration
	buildCheckInterval time.Duration

	jenkins *gojenkins.Jenkins
}

var jenkinsConfigurationFile = flag.String("JenkinsConfigurationFile", "~/.jenkins.yaml", "JenkinsConfigurationFile")

func New() *JenkinsInstance {
	j := &JenkinsInstance{}

	expandedJenkinsConfigurationFile, err := homedir.Expand(*jenkinsConfigurationFile)
	logging.Panic(err)

	yamlhelper.Read(expandedJenkinsConfigurationFile, j)

	// decrypt username and password
	j.username = secrets.Decrypt(j.username)
	j.password = secrets.Decrypt(j.password)

	j.jenkins = gojenkins.CreateJenkins(j.url, j.username, j.password)

	_, err = j.jenkins.Init()
	logging.Panic(err)

	return j
}
