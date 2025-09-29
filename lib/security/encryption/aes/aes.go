package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"os"
)

type AES struct {
	gcm cipher.AEAD
}

func FromFile(path string) (*AES, error) {
	data, err := fileContents(path)
	if err != nil {
		return nil, err
	}

	return New(data)
}

func fileContents(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	return io.ReadAll(f)
}

func New(key []byte) (*AES, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return &AES{gcm: gcm}, nil
}

func (a *AES) Encrypt(data []byte) []byte {
	return a.gcm.Seal(nil, nil, data, nil)
}

func (a *AES) Decrypt(input []byte) ([]byte, error) {
	if len(input) < a.gcm.Overhead() {
		return nil, errors.New("input too short")
	}

	plain, err := a.gcm.Open(nil, nil, input, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
