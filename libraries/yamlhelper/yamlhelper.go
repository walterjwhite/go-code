package yamlhelper

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Read(configurationFile string, out interface{}) {
	log.Debug().Msgf("reading: %v / %v", configurationFile, out)

	yamlFile, err := ioutil.ReadFile(configurationFile)
	logging.Panic(err)

	logging.Panic(yaml.Unmarshal(yamlFile, out))

	log.Debug().Msgf("Read:\n%v", out)
}
