package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/walterjwhite/go-code/lib/security/encryption"
	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize = 16
	pbkdf2Iterations = 310000
	pbkdf2KeyLength = 32
	saltMarker = 0x53414C54 // "SALT" in ASCII
)

type AES struct {
	gcm  cipher.AEAD
	salt []byte // Salt used for PBKDF2 key derivation (nil if key was provided directly)
}

var _ encryption.Encryptor = (*AES)(nil)

func FromFile(path string) (*AES, error) {
	sanitizedPath, err := sanitizeKeyPath(path)
	if err != nil {
		return nil, err
	}

	data, err := fileContents(sanitizedPath)
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

func sanitizeKeyPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("key file path cannot be empty")
	}

	if strings.Contains(path, "..") {
		return "", errors.New("invalid key file path: path traversal not allowed")
	}

	cleanPath := filepath.Clean(path)

	if strings.Contains(cleanPath, "..") {
		return "", errors.New("invalid key file path: resolved path contains traversal")
	}

	return cleanPath, nil
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

func isValidKeyLength(length int) bool {
	return length == 16 || length == 24 || length == 32
}

func NewFromWeakKey(weakKey []byte) (*AES, error) {
	if len(weakKey) == 0 {
		return nil, errors.New("weak key cannot be empty")
	}

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}

	derivedKey := pbkdf2.Key(weakKey, salt, pbkdf2Iterations, pbkdf2KeyLength, sha256.New)

	aesInstance, err := New(derivedKey)
	if err != nil {
		return nil, err
	}
	aesInstance.salt = salt
	return aesInstance, nil
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

	encrypted := a.gcm.Seal(nonce, nonce, data, nil)

	if a.salt != nil {
		result := make([]byte, 4+len(a.salt)+len(encrypted))
		binary.BigEndian.PutUint32(result[:4], saltMarker)
		copy(result[4:4+len(a.salt)], a.salt)
		copy(result[4+len(a.salt):], encrypted)
		return result, nil
	}

	return encrypted, nil
}

func (a *AES) Decrypt(input []byte) ([]byte, error) {
	nonceSize := a.gcm.NonceSize()

	if len(input) >= 4 {
		marker := binary.BigEndian.Uint32(input[:4])
		if marker == saltMarker && a.salt != nil {
			if len(input) < 4+len(a.salt)+nonceSize {
				return nil, errors.New("ciphertext too short: invalid format or corrupted data")
			}
			input = input[4+len(a.salt):]
		}
	}

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
