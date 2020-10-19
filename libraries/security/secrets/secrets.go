package secrets

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/security/encryption"
	"github.com/walterjwhite/go-application/libraries/security/encryption/aes"
	"github.com/walterjwhite/go-application/libraries/security/encryption/providers/file"
	"github.com/walterjwhite/go-application/libraries/security/encryption/providers/ssh"
	"github.com/walterjwhite/go-application/libraries/security/encryption/providers/stdin"

	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/io/yaml"
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

// TODO: this should be moved out to a CLI module?
var (
	secretConfigurationFilePath  = flag.String("c", "~/.config/walterjwhite/secrets.yaml", "SecretsConfigurationFilePath")
	secretFileFlag               = flag.String("f", "", "Secret Key Filename")
	secretProviderFlag           = flag.String("p", "", "Secret Provider")
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

	yaml.Read(filename, SecretsConfigurationInstance)

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
