package property

import (
	"github.com/rs/zerolog/log"
)


func Load(config interface{}) {
	log.Debug().Msgf("before configuration: %v", config)

	LoadFile(config)

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

	log.Debug().Msgf("after configuration: %v", config)
}
