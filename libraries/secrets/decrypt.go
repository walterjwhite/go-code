package secrets

import (
	"log"
	"strings"
)

func Decrypt(secretPath string) string {
	log.Printf("processing secret: %v\n", secretPath)

	initialize()
	setupEncryptionKey()

	data := SecretsConfigurationInstance.encryptionConfiguration.DecryptFile(secretPath)
	return strings.TrimSpace(string(data[:]))
}
