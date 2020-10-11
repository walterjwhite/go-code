package main

import (
	"fmt"

	"github.com/walterjwhite/go-application/libraries/security/secrets"
)

func encrypt() {
	data := readStdIn("secret data (plaintext)")
	fmt.Println(secrets.Base64(secrets.DoEncrypt([]byte(*data))))
}
