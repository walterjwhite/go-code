package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
	"github.com/walterjwhite/go-application/libraries/token/plugins/stdin"
)

var (
	name    = flag.String("name", "", "Secret key name (hierarchy to key, excluding trailing /value, ie. /email/gmail.com/personal/email-address)")
	message = flag.String("message", "", "Commit message")
)

func init() {
	application.Configure()
}

// TODO: add support for flags
// instead of specifying the key type (email, user, pass), use a flag instead (-e, -u, -p)
func main() {
	r := &stdin.StdInReader{PromptMessage: "Please enter secret\n"}
	data := r.Get()
	validatePut(name, message, &data)

	secrets.Encrypt(name, message, []byte(data))
}

func validatePut(name, message, data *string) {
	if len(*name) == 0 {
		logging.Panic(errors.New("No name was provided."))
	}

	if len(*message) == 0 {
		logging.Panic(errors.New("No commit message was provided."))
	}

	if len(*data) == 0 {
		logging.Panic(errors.New("No secret data was provided."))
	}
}
