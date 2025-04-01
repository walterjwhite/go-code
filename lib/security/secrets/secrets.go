package secrets

import (
	"flag"
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/ssh"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/stdin"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/yaml"
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
	secretConfigurationFilePath  = flag.String("secret-conf", "~/.config/walterjwhite/secrets.yaml", "SecretsConfigurationFilePath")
	secretFileFlag               = flag.String("secret-filename", "", "Secret Key Filename")
	secretProviderFlag           = flag.String("secret-provider", "", "Secret Provider")
	SecretsConfigurationInstance *SecretsConfiguration
)

func initialize() {
	if SecretsConfigurationInstance == nil {
		SecretsConfigurationInstance = &SecretsConfiguration{}
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

	initRepository()
}

func initEncryption() {
	SecretsConfigurationInstance.EncryptionConfiguration = &aes.Configuration{Encryption: getEncryption()}
}

func getEncryption() encryption.Encryption {
	switch getSecretProvider(*secretProviderFlag) {
	case File:
		if len(*secretFileFlag) == 0 {
			logging.Panic(fmt.Errorf("expecting secret file to be set"))
		}

		return file.New(*secretFileFlag)
	case Stdin:
		return stdin.New()
	default:
		return ssh.New()
	}
}
