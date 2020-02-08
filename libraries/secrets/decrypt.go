package secrets

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)

	initialize()

	encrypted, err := ioutil.ReadFile(getAbsolute(secretPath))
	logging.Panic(err)

	data := SecretsConfigurationInstance.EncryptionConfiguration.Decrypt(encrypted)
	return strings.TrimSpace(string(data[:]))
}

func getAbsolute(secretPath string) string {
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		return filepath.Join(SecretsConfigurationInstance.RepositoryPath, secretPath, "value")
	}

	return secretPath
}
