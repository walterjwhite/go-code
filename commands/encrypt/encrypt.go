package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/encryption/aes"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/ssh"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	filename          = flag.String("filename", "", "filename to encrypt")
	overwriteExisting = flag.Bool("overwriteExisting", false, "overwriteExisting")
	aesConfiguration  *aes.Configuration
)

func init() {
	aesConfiguration = &aes.Configuration{Encryption: ssh.Instance}

	application.Configure()
}

func main() {
	validateArgumentsFilename()

	data, _ := ioutil.ReadFile(*filename)
	logging.Panic(ioutil.WriteFile(getOutfile(), aesConfiguration.Encrypt(data), 0644))
}

func validateArgumentsFilename() {
	if len(*filename) == 0 {
		logging.Panic(errors.New("Specify a filename"))
	}
}

func getOutfile() string {
	o := strings.Replace(*filename, ".decrypted", ".encrypted", 1)

	if o == *filename && !*overwriteExisting {
		logging.Panic(errors.New("The output filename is the same as the input filename and overwriteExisting is false"))
	}

	return o
}
