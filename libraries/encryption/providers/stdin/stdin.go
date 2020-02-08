package stdin

import (
"bufio"
"github.com/walterjwhite/go-application/libraries/logging"
"fmt"
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

	logging.Panic(fmt.Errorf("No encryption key provided"))
	return nil
}

func (c *Key) GetDecryptionKey() []byte {
	return c.encryptionKey
}

func (c *Key) GetEncryptionKey() []byte {
	return c.encryptionKey
}
