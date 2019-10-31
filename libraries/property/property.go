package property

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/encryption"
)

type Configuration interface {
	HasDefault() bool
	Refreshable() bool
	//Encrypted() bool
	EncryptedFields() []string
}

type ConfigurationReader interface {
	Load(config interface{}, prefix string)
}

var (
	registry []ConfigurationReader
	e        *encryption.EncryptionConfiguration
)

func init() {
	registry = make([]ConfigurationReader, 0)

	registry = append(registry, &defaultConfigurationReader{})
	registry = append(registry, &envConfigurationReader{})
	registry = append(registry, &cliConfigurationReader{})

	e = encryption.New()
}

func Load(config Configuration, prefix string) {
	doLoad(config, prefix)
	handleEncryptedProperties(config)
}

func doLoad(config Configuration, prefix string) {
	for index, value := range registry {
		log.Info().Msgf("reading: %v via %T", index, value)
		value.Load(config, prefix)
	}
}
