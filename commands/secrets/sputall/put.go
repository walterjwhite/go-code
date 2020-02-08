package main

import (
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
	"io/ioutil"
)

func init() {
	application.Configure()
}

func main() {
	secrets.Find(encryptOnMatch, secrets.NewFind()...)
}

func encryptOnMatch(secretPath string) {
	plaintext, err := ioutil.ReadFile(secretPath)
	logging.Panic(err)

	// delete plaintext file (happens when we write back)
	m := "re-encrypt"
	secrets.Encrypt(getName(secretPath), &m, plaintext)
}

func getName(secretPath string) *string {
	n := secretPath[len(secrets.SecretsConfigurationInstance.RepositoryPath)+1:]
	n = n[:len(n)-6]
	return &n
}
