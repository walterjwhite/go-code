package property

import (
	"context"
	"github.com/rs/zerolog/log"
)


type PreLoad interface {
	PreLoad()
}

type PostLoad interface {
	PostLoad(ctx context.Context) error
}

func Load(applicationName string, config interface{}) {
	log.Debug().Msgf("before configuration: %v", config)

	LoadFile(applicationName, config)

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

	log.Debug().Msgf("after configuration: %v", config)
}
