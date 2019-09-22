package main

import (
	"flag"
	"io/ioutil"
	"strings"
	
	"github.com/walterjwhite/go-application/libraries/secrets"
)

var filename = flag.String("filename", "", "filename to decrypt")
var out = flag.String("out", "", "outfile")

func main() {
	_ = application.Configure()
	data := secrets.DecryptFile(getOutfile())
	
	logging.Panic(ioutil.WriteFile(outputFilename, data, 0644))
}

func getOutfile() string {
	var outputFilename string
    if len(*out) > 0 {
       return *out
    }
    
    return strings.Replace(*filename, ".encrypted", ".decrypted", 1)
}
