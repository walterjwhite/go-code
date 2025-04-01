package stdin

import (
	"bufio"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"os"
)

type Key struct {
	encryptionKey []byte
}

func New() *Key {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		keyBytes := scanner.Bytes()
		keyBytes = append(keyBytes, '\n')

		return &Key{encryptionKey: keyBytes}
	}

	logging.Panic(fmt.Errorf("no encryption key provided"))
	return nil
}

func (c *Key) GetDecryptionKey() []byte {
	return c.encryptionKey
}

func (c *Key) GetEncryptionKey() []byte {
	return c.encryptionKey
}
