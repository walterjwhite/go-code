package property

import (
	"github.com/vrischmann/envconfig"
	"github.com/walterjwhite/go-application/libraries/logging"
)

type envConfigurationReader struct{}

func (e *envConfigurationReader) Load(config interface{}, prefix string) {
	logging.Panic(envconfig.InitWithPrefix(config, prefix))
}
