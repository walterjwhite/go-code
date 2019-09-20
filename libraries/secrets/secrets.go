package secrets

import (
	"bufio"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

type SecretsConfiguration struct {
	EncryptionConfiguration encryption.EncryptionConfiguration
	RepositoryRemoteUri     string
	RepositoryPath          string
}

type NoEncryptionKeyProvided struct{}

func (e *NoEncryptionKeyProvided) Error() string {
	return "No key provided"
}

var secretsConfiguration *SecretsConfiguration

// initialize the key
func init() {
	secretsConfiguration = &SecretsConfiguration{}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		secretsConfiguration.EncryptionConfiguration = encryption.EncryptionConfiguration{EncryptionKey: scanner.Bytes()}
	} else {
		logging.Panic(&NoEncryptionKeyProvided{})
	}

	logging.Panic(scanner.Err())
}
