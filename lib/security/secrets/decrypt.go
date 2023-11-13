package secrets

import (
	"encoding/base64"

	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)
	initialize()

	if !filepath.IsAbs(secretPath) {
		secretPath = SecretsConfigurationInstance.RepositoryPath + "/" + secretPath
	}

	encrypted, err := os.ReadFile(secretPath)
	logging.Panic(err)

	data := DoDecrypt(encrypted)
	return strings.TrimSpace(string(data[:]))
}

func DoDecrypt(data []byte) []byte {
	initEncryption()

	return SecretsConfigurationInstance.EncryptionConfiguration.Decrypt(data)
}

func Unbase64(data string) []byte {
	raw, err := base64.StdEncoding.DecodeString(data)
	logging.Panic(err)

	return raw
}
