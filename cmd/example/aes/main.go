package main

import (
	"errors"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		logging.Panic(errors.New("expected exactly 1 argument, the encryption key"))
	}

	aesConf := &aes.Configuration{}
	aesConf.Encryption = file.New(os.Args[1])

	encryptedData := aesConf.Encrypt([]byte("This is a test"))
	fmt.Printf("Encrypted: %s\n", encryptedData)

	decryptedData := aesConf.Decrypt(encryptedData)
	fmt.Printf("Decrypted: %s\n", decryptedData)
}
