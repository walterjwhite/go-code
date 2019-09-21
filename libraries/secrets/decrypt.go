package secrets

import (
	"strings"
)

func Decrypt(secretPath string) string {
	data := SecretsConfigurationInstance.EncryptionConfiguration.DecryptFile(secretPath)
	return strings.TrimSpace(string(data[:]))
}
