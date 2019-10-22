package secrets

import (
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type SecretsConfiguration struct {
	EncryptionConfiguration *encryption.EncryptionConfiguration
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
	if SecretsConfigurationInstance.EncryptionConfiguration != nil {
		return
	}

	SecretsConfigurationInstance.EncryptionConfiguration = encryption.New()
}
