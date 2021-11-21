package main

import (
	"errors"
	"flag"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/secrets"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) < 1 {
		logging.Panic(errors.New("Command is required"))
	}

	switch flag.Args()[0] {
	case "find":
		find()
	case "get":
		get()
	case "put":
		put()
	case "encrypt":
		encrypt()
	case "decrypt":
		decrypt()
	}
}

func onFind(function func(secretPath string), args []string) {
	secrets.Find(function, args...)
}
