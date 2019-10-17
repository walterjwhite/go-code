package secrets

import (
	"github.com/rs/zerolog/log"
	"strings"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)

	initialize()
	setupEncryptionKey()

	data := SecretsConfigurationInstance.encryptionConfiguration.DecryptFile(secretPath)
	return strings.TrimSpace(string(data[:]))
}
