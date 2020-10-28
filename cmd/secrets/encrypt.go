package main

import (
	"fmt"

	"github.com/walterjwhite/go/lib/security/secrets"
)

func encrypt() {
	data := readStdIn("secret data (plaintext)")
	fmt.Println(secrets.Base64(secrets.DoEncrypt([]byte(*data))))
}
