package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"

	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/encryption"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"github.com/walterjwhite/go-application/libraries/secrets"
)

var filename = flag.String("filename", "", "filename to decrypt")
var overwriteExisting = flag.Bool("overwriteExisting", false, "overwriteExisting")

func main() {
	_ = application.Configure()

	validateArgumentsFilename()

	e := encryption.New()
	data := e.DecryptFile(*filename)

	logging.Panic(ioutil.WriteFile(getOutfile(), data, 0644))
}

func validateArgumentsFilename() {
	if len(*filename) == 0 {
		logging.Panic(errors.New("Specify a filename"))
	}
}

func getOutfile() string {
	o := strings.Replace(*filename, ".encrypted", ".decrypted", 1)

	if o == *filename && !*overwriteExisting {
		logging.Panic(errors.New("The output filename is the same as the input filename and overwriteExisting is false"))
	}

	return o
}
