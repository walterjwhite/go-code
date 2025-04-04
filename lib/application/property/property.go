package property

import (
	"github.com/rs/zerolog/log"
)


func Load(config interface{}) {
	LoadFile(config)

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

	log.Debug().Msgf("configuration: %v", config)
}
