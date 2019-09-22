package main

import (
	"flag"
	"io/ioutil"
	"strings"
	
	"github.com/walterjwhite/go-application/libraries/secrets"
)

var filename = flag.String("filename", "", "filename to encrypt")

func main() {
	_ = application.Configure()
	
	data, _ := ioutil.ReadFile(*filename)
	
	secrets.EncryptFile(getOutfile(), data)
}

func getOutfile() string {
	return strings.Replace(*filename, ".decrypted", ".encrypted", 1)
}
