package secrets

import (
	"bufio"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"os"
)

type SecretsConfiguration struct {
	EncryptionConfiguration encryption.EncryptionConfiguration
	RepositoryRemoteUri     string
	RepositoryPath          string
}

type NoEncryptionKeyProvided struct{}

const secretConfigurationFilePath = "~/.secrets.yaml"

func (e *NoEncryptionKeyProvided) Error() string {
	return "No key provided"
}

var secretsConfiguration *SecretsConfiguration

// initialize the key
func init() {
	secretsConfiguration = &SecretsConfiguration{}

	filename, err := homedir.Expand(secretConfigurationFilePath)
	logging.Panic(err)

	yamlhelper.Read(filename, secretsConfiguration)

	secretsConfiguration.RepositoryPath, err = homedir.Expand(secretsConfiguration.RepositoryPath)
	logging.Panic(err)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		secretsConfiguration.EncryptionConfiguration = encryption.EncryptionConfiguration{EncryptionKey: scanner.Bytes()}
	} else {
		logging.Panic(&NoEncryptionKeyProvided{})
	}

	logging.Panic(scanner.Err())
}
