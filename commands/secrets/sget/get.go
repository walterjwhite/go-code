package main

import (
	"flag"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/secrets"
)

var isDisplayOnStdOut = flag.Bool("StdOut", false, "display secret on StdOut")

// TODO: add support for flags
// instead of specifying the key type (email, user, pass), use a flag instead (-e, -u, -p)
func init() {
	application.Configure()
}

func main() {
	secrets.Find(decryptOnMatch, secrets.NewFind()...)
}

func decryptOnMatch(secretPath string) {
	secretText := secrets.Decrypt(secretPath)

	if *isDisplayOnStdOut {
		fmt.Println(secretText)
	} else {
		logging.Panic(clipboard.WriteAll(secretText))
	}
}
