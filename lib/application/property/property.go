package property

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)


type PreLoad interface {
	PreLoad()
}

type PostLoad interface {
	PostLoad(ctx context.Context) error
}

func Load(applicationName string, config interface{}) {
	log.Debug().Msgf("before configuration: %v", config)

	logging.Warn(LoadFile(applicationName, config), "LoadFile")

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

	log.Debug().Msgf("after configuration: %v", config)
}
