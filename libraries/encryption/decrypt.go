package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (c *EncryptionConfiguration) Decrypt(data []byte) []byte {
	block, err := aes.NewCipher(c.EncryptionKey)
	logging.Panic(err)

	gcm, err := cipher.NewGCM(block)
	logging.Panic(err)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	logging.Panic(err)

	return plaintext
}
