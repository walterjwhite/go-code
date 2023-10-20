package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"

	"os"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/secrets"
)

var (
	getFlagSet = flag.NewFlagSet("get", flag.ExitOnError)

	getOutputTarget = getFlagSet.String("o", "c", "display secret on (c=>Clipboard, s=>StdOut, f=>file)")
	removeSpaces    = getFlagSet.Bool("s", false, "remove spaces")
)

func decryptOnMatch(secretPath string) {
	secretText := secrets.Decrypt(secretPath)

	if *removeSpaces {
		secretText = strings.ReplaceAll(secretText, " ", "")
	}

	switch *getOutputTarget {
	case "s":
		fmt.Println(secretText)
	case "f":
		logging.Panic(os.WriteFile(secretPath+".dec", []byte(secretText), 0644))
	default:
		logging.Panic(clipboard.WriteAll(secretText))
	}
}

func get() {
	logging.Panic(getFlagSet.Parse(flag.Args()[1:]))
	if len(getFlagSet.Args()) == 1 {
		decryptOnMatch(filepath.Join(getFlagSet.Args()[0], "value"))
		return
	}

	onFind(decryptOnMatch, getFlagSet.Args())
}
