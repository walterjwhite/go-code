package aes

import (
	aesl "crypto/aes"
	"crypto/cipher"
	"github.com/walterjwhite/go/lib/application/logging"
)

func (c *Configuration) Decrypt(data []byte) []byte {
	block, err := aesl.NewCipher(c.Encryption.GetDecryptionKey())
	logging.Panic(err)

	gcm, err := cipher.NewGCM(block)
	logging.Panic(err)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	logging.Panic(err)

	return plaintext
}
