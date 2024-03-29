package file

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

type Key struct {
	encryptionKey []byte
}

func New(filename string) *Key {
	data, err := os.ReadFile(filename)
	logging.Panic(err)

	return &Key{encryptionKey: data}
}

func (c *Key) GetDecryptionKey() []byte {
	return c.encryptionKey
}

func (c *Key) GetEncryptionKey() []byte {
	return c.encryptionKey
}
