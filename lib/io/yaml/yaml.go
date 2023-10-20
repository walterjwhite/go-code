package yaml

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"gopkg.in/yaml.v2"
	"os"
)

func Read(configurationFile string, out interface{}) {
	log.Debug().Msgf("reading: %v / %v", configurationFile, out)

	yamlFile, err := os.ReadFile(configurationFile)
	logging.Panic(err)

	logging.Panic(yaml.Unmarshal(yamlFile, out))

	log.Debug().Msgf("Read:\n%v", out)
}

func Write(in interface{}, outFile string) {
	buf, err := yaml.Marshal(in)
	logging.Panic(err)

	logging.Panic(os.WriteFile(outFile, buf /*os.ModePerm*/, 0644))

	log.Debug().Msgf("Wrote:\n%v", outFile)
}
