package main

import (
	"github.com/walterjwhite/go-application/libraries/application"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/encryption/aes"
	"github.com/walterjwhite/go-application/libraries/encryption/providers/ssh"
)

var (
	e *aes.Configuration
)

func init() {
	application.Configure()

	e = &aes.Configuration{Encryption: ssh.Instance}
}

func main() {
	//ssh.List()
	//ssh.GetDecryptionKey()
	//ssh.GetEncryptionKey()

	// encrypt something
	encrypted := e.Encrypt([]byte("Something"))
	log.Info().Msgf("encrypted: %v", string(encrypted))

	// decrypt something
	decrypted := e.Decrypt(encrypted)

	log.Info().Msgf("decrypted: %v", string(decrypted))
}
