package encryption

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

func (c *EncryptionConfiguration) EncryptFile(filename string, data []byte) {
	f, err := os.Create(filename)
	logging.Panic(err)

	defer f.Close()

	_, err = f.Write(c.Encrypt(data))
	logging.Panic(err)
}
