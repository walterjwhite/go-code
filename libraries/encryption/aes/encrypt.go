package aes

import (
	aesl "crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
)

func (c *Configuration) Encrypt(data []byte) []byte {
	log.Warn().Msgf("key size: %v", len(c.Encryption.GetEncryptionKey()))

	block, err := aesl.NewCipher(c.Encryption.GetEncryptionKey())
	logging.Panic(err)

	gcm, err := cipher.NewGCM(block)
	logging.Panic(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	logging.Panic(err)

	return gcm.Seal(nonce, nonce, data, nil)
}

/*
func generateSalt(n uint32) []byte {
    b := make([]byte, n)
    _, err := rand.Read(b)
    logging.Panic(err)

    return b
}

//key := argon2.Key([]byte("some password"), salt, 3, 32*1024, 4, 32)
*/
