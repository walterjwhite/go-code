package server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

const (
	systemBuildConfigurationFile = "/etc/system"
)

type SystemBuildConfiguration struct {
	BuildDate string
}

var systemBuildConfiguration SystemBuildConfiguration

// initialize once on start
var buildDateTime string

func (s *Server) BuildDateTime(args *Args, response *string) error {
	*response = buildDateTime

	log.Printf("response: %v\n", response)
	log.Printf("buildDateTime: %v\n", buildDateTime)
	return nil
}

func RefreshBuildDateTime() {
	// TODO: make this configurable ...
	yamlHelper.Read(systemBuildConfigurationFile, &systemBuildConfiguration)

	buildDateTime = "Built-On: " + systemBuildConfiguration.BuildDate
	log.Printf("Build Date/Time: %v", buildDateTime)
}
