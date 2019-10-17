package secrets

import (
	"bufio"
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"os"
)

type SecretsConfiguration struct {
	encryptionConfiguration encryption.EncryptionConfiguration
	RepositoryRemoteUri     string
	RepositoryPath          string
}

type NoEncryptionKeyProvided struct{}

var secretConfigurationFilePath = flag.String("SecretsConfigurationFilePath", "~/.secrets.yaml", "SecretsConfigurationFilePath")

func (e *NoEncryptionKeyProvided) Error() string {
	return "No key provided"
}

var SecretsConfigurationInstance *SecretsConfiguration

// initialize the key
func initialize() {
	if SecretsConfigurationInstance != nil {
		return
	}

	SecretsConfigurationInstance = &SecretsConfiguration{}

	filename, err := homedir.Expand(*secretConfigurationFilePath)
	logging.Panic(err)

	yamlhelper.Read(filename, SecretsConfigurationInstance)

	translatedRepositoryPath, err := homedir.Expand(SecretsConfigurationInstance.RepositoryPath)
	SecretsConfigurationInstance.RepositoryPath = translatedRepositoryPath
	logging.Panic(err)

	setupRepository()
}

func setupEncryptionKey() {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		keyBytes := scanner.Bytes()
		keyBytes = append(keyBytes, '\n')

		SecretsConfigurationInstance.encryptionConfiguration = encryption.EncryptionConfiguration{EncryptionKey: keyBytes}
	} else {
		logging.Panic(&NoEncryptionKeyProvided{})
	}

	logging.Panic(scanner.Err())
}
