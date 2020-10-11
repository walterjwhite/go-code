package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"github.com/walterjwhite/go-application/libraries/security/secrets"
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
		onFind(printOnMatch, flag.Args()[1:])
	case "get":
		logging.Panic(getFlagSet.Parse(flag.Args()[1:]))
		onFind(decryptOnMatch, getFlagSet.Args())
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
