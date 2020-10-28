package main

import (
	"flag"
	"fmt"

	"github.com/atotto/clipboard"

	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/security/secrets"
	"io/ioutil"
)

var (
	getFlagSet = flag.NewFlagSet("get", flag.ExitOnError)

	getOutputTarget = getFlagSet.String("o", "c", "display secret on (c=>Clipboard, s=>StdOut, f=>file)")
)

func decryptOnMatch(secretPath string) {
	secretText := secrets.Decrypt(secretPath)

	switch *getOutputTarget {
	case "s":
		fmt.Println(secretText)
	case "f":
		logging.Panic(ioutil.WriteFile(secretPath+".dec", []byte(secretText), 0644))
	default:
		logging.Panic(clipboard.WriteAll(secretText))
	}
}
