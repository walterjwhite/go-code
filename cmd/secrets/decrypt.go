package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/security/secrets"
)

var (
	decryptFlagSet = flag.NewFlagSet("decrypt", flag.ExitOnError)

	decryptStdOut = putFlagSet.Bool("o", false, "print decrypted value to StdOut (defaults to clipboard)")
)

func decrypt() {
	logging.Panic(decryptFlagSet.Parse(flag.Args()[1:]))

	data := readStdIn("secret data (cipher text)")

	secretText := string(secrets.DoDecrypt(secrets.Unbase64(*data)))[:]

	if *decryptStdOut {
		fmt.Println(secretText)
	} else {
		logging.Panic(clipboard.WriteAll(secretText))
	}
}
