package property

import (
	"github.com/rs/zerolog/log"
)

type ConfigurationReader interface {
	Load(config interface{}, prefix string)
}

var (
	readerRegistry []ConfigurationReader
)

func init() {
	readerRegistry = make([]ConfigurationReader, 0)

	readerRegistry = append(readerRegistry, &defaultConfigurationReader{})
	readerRegistry = append(readerRegistry, &envConfigurationReader{})
	readerRegistry = append(readerRegistry, &cliConfigurationReader{})
	readerRegistry = append(readerRegistry, &cliConfigurationReader{})

	readerRegistry = append(readerRegistry, &encryptionReader{})
}

func Load(config interface{}, prefix string) {
	for index, value := range readerRegistry {
		log.Info().Msgf("reading: %v via %T", index, value)
		value.Load(config, prefix)
	}
}
