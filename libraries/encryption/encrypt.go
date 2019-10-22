package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
)

func (c *EncryptionConfiguration) Encrypt(data []byte) []byte {
	block, err := aes.NewCipher(c.encryptionKey)
	logging.Panic(err)

	gcm, err := cipher.NewGCM(block)
	logging.Panic(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	logging.Panic(err)

	return gcm.Seal(nonce, nonce, data, nil)
}
