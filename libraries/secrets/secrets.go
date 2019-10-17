package secrets

import (
	"bufio"
	"errors"
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"os"
)

type SecretsConfiguration struct {
	encryptionConfiguration *encryption.EncryptionConfiguration
	RepositoryRemoteUri     string
	RepositoryPath          string
}

var secretConfigurationFilePath = flag.String("SecretsConfigurationFilePath", "~/.secrets.yaml", "SecretsConfigurationFilePath")

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
	if SecretsConfigurationInstance.encryptionConfiguration != nil {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		keyBytes := scanner.Bytes()
		keyBytes = append(keyBytes, '\n')

		SecretsConfigurationInstance.encryptionConfiguration = &encryption.EncryptionConfiguration{EncryptionKey: keyBytes}
	} else {
		logging.Panic(errors.New("No encryption key provided"))
	}

	logging.Panic(scanner.Err())
}
