package maven

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf(".profiles/%v.yaml", profile))
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
