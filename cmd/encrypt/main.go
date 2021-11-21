package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/ssh"
)

var (
	filename          = flag.String("f", "", "filename to encrypt/decrypt")
	overwriteExisting = flag.Bool("o", false, "overwrite existing")
	isDecrypt         = flag.Bool("d", false, "(encrypt by default) -d to decrypt")

	aesConfiguration *aes.Configuration
)

func init() {
	aesConfiguration = &aes.Configuration{Encryption: ssh.Instance}

	application.Configure()
}

func main() {
	if len(*filename) == 0 {
		logging.Panic(errors.New("Specify a filename"))
	}

	logging.Panic(ioutil.WriteFile(getOutfile(), doWork(), 0644))
}

func getOutfile() string {
	var o string
	if *isDecrypt {
		o = strings.Replace(*filename, ".encrypted", ".decrypted", 1)
	} else {
		o = strings.Replace(*filename, ".decrypted", ".encrypted", 1)
	}

	if o == *filename && !*overwriteExisting {
		logging.Panic(errors.New("The output filename is the same as the input filename and overwriteExisting is false"))
	}

	return o
}

func doWork() []byte {
	data, _ := ioutil.ReadFile(*filename)
	if *isDecrypt {
		return aesConfiguration.Decrypt(data)
	}

	return aesConfiguration.Encrypt(data)
}
