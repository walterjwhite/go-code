package property

import (
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
)

func (c *Configuration) LoadEnv(config interface{}) {
	if len(c.Path) > 0 {
		err := envconfig.InitWithPrefix(config, c.Path)
		log.Warn().Msgf("Error reading properties from env: %v", err)
	} else {
		err := envconfig.Init(config)
		log.Warn().Msgf("Error reading properties from env: %v", err)
	}
}
