package secrets

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Decrypt(secretPath string) string {
	log.Debug().Msgf("processing secret: %v", secretPath)
	initialize()

	if !filepath.IsAbs(secretPath) {
		secretPath = SecretsConfigurationInstance.RepositoryPath + "/" + secretPath
	}

	encrypted, err := ioutil.ReadFile(secretPath)
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
