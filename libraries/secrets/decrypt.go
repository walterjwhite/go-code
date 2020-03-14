package secrets

import (
	"encoding/base64"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)

	if !filepath.IsAbs(secretPath) {
		initialize()

		secretPath = SecretsConfigurationInstance.RepositoryPath + "/" + secretPath
	}

	encrypted, err := ioutil.ReadFile(secretPath)
	logging.Panic(err)

	data := DoDecrypt(encrypted)
	return strings.TrimSpace(string(data[:]))
}

func DoDecrypt(data []byte) []byte {
	initialize()

	return SecretsConfigurationInstance.EncryptionConfiguration.Decrypt(data)
}

func Unbase64(data string) []byte {
	raw, err := base64.StdEncoding.DecodeString(data)
	logging.Panic(err)

	return raw
}
