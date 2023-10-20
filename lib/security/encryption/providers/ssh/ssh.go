package ssh

import (
	"crypto/sha256"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func New() *Conf {
	key, err := initAgent()
	if err == nil {
		return &Conf{key: getUsableKey(key)}
	}

	log.Warn().Msgf("error initializing agent: %v", err)

	key, err = initDirect()
	if err == nil {
		return &Conf{key: getUsableKey(key)}
	}

	logging.Panic(err)
	return nil
}

func getUsableKey(data []byte) []byte {
	k := sha256.Sum256(data)
	return k[:]
}
