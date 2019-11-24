package main

import (
	"errors"
	"flag"
	"io/ioutil"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
)

var (
	name    = flag.String("name", "", "Secret key name (hierarchy to key, excluding trailing /value, ie. /email/gmail.com/personal/email-address)")
	message = flag.String("message", "", "Commit message")
	source  = flag.String("source", "", "source file")
)

func init() {
	application.Configure()
}

// TODO: add support for flags
// instead of specifying the key type (email, user, pass), use a flag instead (-e, -u, -p)
func main() {
	validatePut(name, message, source)

	data, err := ioutil.ReadFile(*source)
	logging.Panic(err)

	secrets.Encrypt(name, message, data)
}

func validatePut(name *string, message *string, source *string) {
	if len(*name) == 0 {
		logging.Panic(errors.New("No name was provided."))
	}

	if len(*message) == 0 {
		logging.Panic(errors.New("No commit message was provided."))
	}

	if len(*source) == 0 {
		logging.Panic(errors.New("No secret data was provided."))
	}
}
