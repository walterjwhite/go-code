package run

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type Application struct {
	Name  string
	Command string
	Arguments []string
	LogMatcher string
	Environment  []string
}

type Configuration struct {
	Applications []Application
}

func (a *Application) getConf(application string) *Application {
	yamlhelper.Read(fmt.Sprintf(".applications/%v.yaml", application), a)

	return c
}
