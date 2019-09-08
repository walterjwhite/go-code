package yamlhelper

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func Read(configurationFile string, out interface{}) {
	log.Printf("reading: %v / %v\n", configurationFile, out)
	yamlFile, err := ioutil.ReadFile(configurationFile)
	logging.Panic(err)

	logging.Panic(yaml.Unmarshal(yamlFile, out))
	// TODO: should this be a separate exception?
	/*
		if err != nil {
			log.Fatalf("Error unmarshalling file: %v / %v / %v\n", configurationFile, yamlFile, err)
		}
	*/

	log.Printf("Read:\n%v\n", out)
}
