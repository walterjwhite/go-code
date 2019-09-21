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

var SecretsConfigurationInstance *SecretsConfiguration

// initialize the key
func init() {
	SecretsConfigurationInstance = &SecretsConfiguration{}

	filename, err := homedir.Expand(secretConfigurationFilePath)
	logging.Panic(err)

	yamlhelper.Read(filename, SecretsConfigurationInstance)

	translatedRepositoryPath, err := homedir.Expand(SecretsConfigurationInstance.RepositoryPath)
	SecretsConfigurationInstance.RepositoryPath = translatedRepositoryPath
	logging.Panic(err)

	scanner := bufio.NewScanner(os.Stdin)
	keyBytes := scanner.Bytes()
	keyBytes = append(keyBytes, '\n')
	if scanner.Scan() {
		SecretsConfigurationInstance.EncryptionConfiguration = encryption.EncryptionConfiguration{EncryptionKey: keyBytes}
	} else {
		logging.Panic(&NoEncryptionKeyProvided{})
	}

	logging.Panic(scanner.Err())
}
