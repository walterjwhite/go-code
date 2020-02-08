package aes

import (
	aesl "crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
)

func (c *Configuration) Encrypt(data []byte) []byte {
	block, err := aesl.NewCipher(c.Encryption.GetEncryptionKey())
	logging.Panic(err)

	gcm, err := cipher.NewGCM(block)
	logging.Panic(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	logging.Panic(err)

	return gcm.Seal(nonce, nonce, data, nil)
}
