package secrets

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/encryption/aes"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/file"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/ssh"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/stdin"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type SecretsConfiguration struct {
	EncryptionConfiguration *aes.Configuration
	RepositoryRemoteUri     string
	RepositoryPath          string
}

type SecretProvider int

const (
	SSH SecretProvider = iota
	Stdin
	File
)

func (p SecretProvider) String() string {
	return [...]string{"SSH", "Stdin", "File"}[p]
}

func getSecretProvider(providerName string) SecretProvider {
	switch providerName {
	case "File":
		return File
	case "Stdin":
		return Stdin
	default:
		return SSH
	}
}

var (
	secretConfigurationFilePath  = flag.String("SecretsConfigurationFilePath", "~/.secrets.yaml", "SecretsConfigurationFilePath")
	secretFileFlag               = flag.String("SecretKey", "", "Secret Key Filename")
	secretProviderFlag           = flag.String("SecretProvider", "", "Secret Provider")
	SecretsConfigurationInstance *SecretsConfiguration
)

// initialize the key
func initialize() {
	if SecretsConfigurationInstance == nil {
		SecretsConfigurationInstance = &SecretsConfiguration{EncryptionConfiguration: &aes.Configuration{Encryption: getEncryption()}}
	} else {
		if len(SecretsConfigurationInstance.RepositoryPath) > 0 {
			return
		}
	}

	filename, err := homedir.Expand(*secretConfigurationFilePath)
	logging.Panic(err)

	yamlhelper.Read(filename, SecretsConfigurationInstance)

	translatedRepositoryPath, err := homedir.Expand(SecretsConfigurationInstance.RepositoryPath)
	SecretsConfigurationInstance.RepositoryPath = translatedRepositoryPath
	logging.Panic(err)

	setupRepository()
}

func getEncryption() encryption.Encryption {
	switch getSecretProvider(*secretProviderFlag) {
	case File:
		if len(*secretFileFlag) == 0 {
			logging.Panic(fmt.Errorf("Expecting secret file to be set"))
		}

		return file.New(*secretFileFlag)
	case Stdin:
		return stdin.New()
	default:
		return ssh.Instance
	}
}
