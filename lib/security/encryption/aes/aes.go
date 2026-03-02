package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/walterjwhite/go-code/lib/security/encryption"
	"golang.org/x/crypto/pbkdf2"
)

type AES struct {
	gcm cipher.AEAD
}

var _ encryption.Encryptor = (*AES)(nil)

func FromFile(path string) (*AES, error) {
	data, err := fileContents(path)
	if err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(data))
	if err == nil && isValidKeyLength(len(key)) {
		return New(key)
	}

	if isValidKeyLength(len(data)) {
		return New(data)
	}

	return NewFromWeakKey(data)
}

func FromEnv(envVarName string) (*AES, error) {
	keyStr := os.Getenv(envVarName)
	if keyStr == "" {
		return nil, errors.New("encryption key environment variable not set: " + envVarName)
	}

	key, err := hex.DecodeString(keyStr)
	if err == nil && isValidKeyLength(len(key)) {
		return New(key)
	}

	return NewFromWeakKey([]byte(keyStr))
}

func validateKeyLength(key []byte) ([]byte, error) {
	if !isValidKeyLength(len(key)) {
		return nil, errors.New("invalid key length: must be 16, 24, or 32 bytes")
	}
	return key, nil
}

func isValidKeyLength(length int) bool {
	return length == 16 || length == 24 || length == 32
}

func NewFromWeakKey(weakKey []byte) (*AES, error) {
	if len(weakKey) == 0 {
		return nil, errors.New("weak key cannot be empty")
	}

	const iterations = 100000
	const keyLength = 32

	h := sha256.Sum256(weakKey)
	salt := h[:]

	derivedKey := pbkdf2.Key(weakKey, salt, iterations, keyLength, sha256.New)
	return New(derivedKey)
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
	if len(key) == 0 {
		return nil, errors.New("encryption key cannot be empty")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCM mode: %w", err)
	}

	return &AES{gcm: gcm}, nil
}

func (a *AES) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate random nonce: %w", err)
	}
	if len(nonce) != a.gcm.NonceSize() {
		return nil, errors.New("nonce generation failed: invalid size")
	}
	return a.gcm.Seal(nonce, nonce, data, nil), nil
}

func (a *AES) Decrypt(input []byte) ([]byte, error) {
	nonceSize := a.gcm.NonceSize()
	if len(input) < nonceSize {
		return nil, errors.New("ciphertext too short: invalid format or corrupted data")
	}

	if len(input) <= nonceSize {
		return nil, errors.New("ciphertext too short: insufficient data for decryption")
	}

	nonce, ciphertext := input[:nonceSize], input[nonceSize:]

	if len(nonce) != nonceSize {
		return nil, errors.New("invalid nonce size")
	}

	plain, err := a.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid ciphertext or wrong key")
	}

	return plain, nil
}
