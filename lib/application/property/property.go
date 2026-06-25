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

func Load(applicationName string, config any) {
	log.Debug().Msgf("loading configuration for: %T", config)

	logging.Warn(LoadFile(applicationName, config), "LoadFile")

	LoadEnv(config)
	LoadCli(config)
	LoadSecrets(config)

	validateRequiredFields(config)

	log.Debug().Msgf("configuration loaded for: %T", config)
}
