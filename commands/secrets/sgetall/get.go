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
	secrets.Find(decryptOnMatch, secrets.NewFind()...)
}

func decryptOnMatch(secretPath string) {
	secretText := secrets.Decrypt(secretPath)

	logging.Panic(ioutil.WriteFile(secretPath+".dec", []byte(secretText), 0644))
}
