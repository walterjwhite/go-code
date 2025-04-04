package property

import (
	"github.com/rs/zerolog/log"
	"github.com/vrischmann/envconfig"
)

func LoadEnv(config interface{}) {
	err := envconfig.Init(config)
	log.Warn().Msgf("Error reading properties from env: %v", err)
}
