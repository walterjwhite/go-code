package aes

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
	"os"
)

func (c *Configuration) EncryptFile(inFilename, outFilename string) {
	data, err := ioutil.ReadFile(inFilename)
	logging.Panic(err)

	f, err := os.Create(outFilename)
	logging.Panic(err)

	defer f.Close()

	_, err = f.Write(c.Encrypt(data))
	logging.Panic(err)
}
