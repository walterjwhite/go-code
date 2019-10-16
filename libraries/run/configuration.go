package run

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"path/filepath"
)

type Application struct {
	Name        string
	Command     string
	Arguments   []string
	LogMatcher  string
	Environment []string
}

type Configuration struct {
	Applications []Application
}

var applicationsPath = flag.String("RunApplicationsPath", "~/.applications", "RunApplicationsPath")

func (a *Application) getConf(application string) *Application {
	f := filepath.Join(*applicationsPath, fmt.Sprintf("%v.yaml", application))
	f, err := homedir.Expand(f)
	logging.Panic(err)

	yamlhelper.Read(f, a)

	return a
}
