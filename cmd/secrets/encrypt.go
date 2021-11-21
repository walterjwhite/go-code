package main

import (
	"fmt"

	"github.com/walterjwhite/go-code/lib/security/secrets"
)

func encrypt() {
	data := readStdIn("secret data (plaintext)")
	fmt.Println(secrets.Base64(secrets.DoEncrypt([]byte(*data))))
}
