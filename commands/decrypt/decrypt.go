package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/encryption/aes"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/ssh"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"github.com/walterjwhite/go-application/libraries/secrets"
)

var (
	filename = flag.String("filename", "", "filename to decrypt")
	//output            = flag.String("output", "", "output filename")
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
	logging.Panic(ioutil.WriteFile(getOutfile(), aesConfiguration.Decrypt(data), 0644))
}

func validateArgumentsFilename() {
	if len(*filename) == 0 {
		logging.Panic(errors.New("Specify a filename"))
	}
}

func getOutfile() string {
	o := doGetOutfile()

	log.Info().Msgf("Writing decrypted contents to: %v", o)

	if o == *filename && !*overwriteExisting {
		logging.Panic(errors.New("The output filename is the same as the input filename and overwriteExisting is false"))
	}

	log.Info().Msgf("Writing decrypted contents to: %v", o)
	return o
}

func doGetOutfile() string {
	/*if len(*output) > 0 {
		return *output
	}*/

	return strings.Replace(*filename, ".encrypted", ".decrypted", 1)
}
