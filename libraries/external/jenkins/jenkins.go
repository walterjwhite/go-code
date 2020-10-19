package jenkins

import (
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"gopkg.in/bndr/gojenkins.v1"
	"time"
)

func (c *Instance) SecretFields() []string {
	return []string{"Username", "Password"}
}

type Instance struct {
	Url string

	Username string
	Password string

	BuildTimeout       *time.Duration
	BuildCheckInterval *time.Duration

	jenkins *gojenkins.Jenkins
}

func (i *Instance) setup() {
	if i.jenkins != nil {
		return
	}

	i.jenkins = gojenkins.CreateJenkins(nil, i.Url, i.Username, i.Password)

	_, err := i.jenkins.Init()
	logging.Panic(err)
}
