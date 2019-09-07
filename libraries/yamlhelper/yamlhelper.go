package yamlhelper

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func Read(configurationFile string, out interface{}) {
	log.Printf("reading: %v / %v\n", configurationFile, out)
	yamlFile, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		log.Fatalf("Error unmarshalling file: %v / %v / %v\n", configurationFile, yamlFile, err)
	}

	log.Printf("Read:\n%v\n", out)
}
