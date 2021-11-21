package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/secrets"
	"github.com/walterjwhite/go-code/lib/utils/token/plugins/stdin"
)

var (
	putFlagSet = flag.NewFlagSet("put", flag.ExitOnError)

	putKeyPath = putFlagSet.String("key", "", "key path ie. web/google/email-address")
	putMessage = putFlagSet.String("message", "", "commit message")
)

// TODO: add support for flags (path, email address, username, password)
func put() {
	logging.Panic(putFlagSet.Parse(flag.Args()[1:]))

	if len(*putKeyPath) == 0 {
		logging.Panic(errors.New("key is required."))
	}

	if len(*putMessage) == 0 {
		logging.Panic(errors.New("message is required."))
	}

	data := readStdIn("secret data")

	secrets.Encrypt(putKeyPath, putMessage, []byte(*data))
}

func readStdIn(message string) *string {
	r := &stdin.StdInReader{PromptMessage: fmt.Sprintf("Please enter %s\n", message)}
	data := r.Get()

	if len(data) == 0 {
		logging.Panic(fmt.Errorf("No %s was provided.", data))
	}

	return &data
}
