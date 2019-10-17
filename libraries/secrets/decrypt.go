package secrets

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)

	initialize()
	setupEncryptionKey()

	data := SecretsConfigurationInstance.encryptionConfiguration.DecryptFile(getAbsolute(secretPath))
	return strings.TrimSpace(string(data[:]))
}

func getAbsolute(secretPath string) string {
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		return filepath.Join(SecretsConfigurationInstance.RepositoryPath, secretPath, "value")
	}

	return secretPath
}
