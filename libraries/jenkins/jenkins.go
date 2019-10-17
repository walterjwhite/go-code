package jenkins

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"gopkg.in/bndr/gojenkins.v1"
	"path/filepath"
	"time"
)

type JenkinsInstance struct {
	Url      string
	Username string
	Password string

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
	j.Username = secrets.Decrypt(filepath.Join(secrets.SecretsConfigurationInstance.RepositoryPath, j.Username))
	j.Password = secrets.Decrypt(filepath.Join(secrets.SecretsConfigurationInstance.RepositoryPath, j.Password))

	j.jenkins = gojenkins.CreateJenkins(j.Url, j.Username, j.Password)

	_, err = j.jenkins.Init()
	logging.Panic(err)

	return j
}
