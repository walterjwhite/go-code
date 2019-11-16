package property

import (
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
)

type envConfigurationReader struct{}

func (e *envConfigurationReader) Load(config interface{}, prefix string) {
	if len(prefix) > 0 {
		err := envconfig.InitWithPrefix(config, prefix)
		log.Warn().Msgf("Error reading properties from env: %v", err)
	} else {
		err := envconfig.Init(config)
		log.Warn().Msgf("Error reading properties from env: %v", err)
	}
}
