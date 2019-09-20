package encryption

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
)

func (c *EncryptionConfiguration) DecryptFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	logging.Panic(err)

	return c.Decrypt(data)
}
