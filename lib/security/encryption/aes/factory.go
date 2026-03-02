package aes

import (
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

func NewAESFromFile(path string) (encryption.Encryptor, error) {
	return FromFile(path)
}

func NewAES(key []byte) (encryption.Encryptor, error) {
	return New(key)
}

func NewAESFromEnv(envVarName string) (encryption.Encryptor, error) {
	return FromEnv(envVarName)
}
