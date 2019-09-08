package run

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type Application struct {
	Index         int
	Configuration Configuration
}

type Configuration struct {
	Applications []string
	Environment  []string
	Jvm          []string

	DebugPorts []int
}

const DebugArguments = "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=%d"
const DebugPortStart = 5005

func (c *Configuration) getConf(profile string) *Configuration {
	yamlhelper.Read(fmt.Sprintf(".profiles/%v.yaml", profile), c)

	return c
}
