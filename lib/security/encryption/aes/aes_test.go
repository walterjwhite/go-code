package aes

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/walterjwhite/go-code/lib/security/encryption"
)

func generateSecureKey() []byte {
	key := make([]byte, 32) // 32 bytes = AES-256
	if _, err := rand.Read(key); err != nil {
		panic("failed to generate secure random key: " + err.Error())
	}
	return key
}

func TestAESEncryptor(t *testing.T) {
	key := generateSecureKey()

	t.Run("Implements Encryptor interface", func(t *testing.T) {
		var _ encryption.Encryptor = (*AES)(nil)
	})

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		encryptor, err := New(key)
		if err != nil {
			t.Fatalf("Failed to create encryptor: %v", err)
		}

		original := []byte("This is secret data")

		encrypted, err := encryptor.Encrypt(original)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		if bytes.Equal(original, encrypted) {
			t.Error("Encrypted data is the same as original")
		}

		decrypted, err := encryptor.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt failed: %v", err)
		}

		if !bytes.Equal(original, decrypted) {
			t.Errorf("Decrypted data doesn't match original.\nOriginal: %s\nDecrypted: %s", original, decrypted)
		}
	})

	t.Run("Encrypt empty data", func(t *testing.T) {
		encryptor, err := New(key)
		if err != nil {
			t.Fatalf("Failed to create encryptor: %v", err)
		}

		original := []byte{}

		encrypted, err := encryptor.Encrypt(original)
		if err != nil {
			t.Fatalf("Encrypt empty data failed: %v", err)
		}

		decrypted, err := encryptor.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt empty data failed: %v", err)
		}

		if !bytes.Equal(original, decrypted) {
			t.Error("Empty data round-trip failed")
		}
	})

	t.Run("Decrypt invalid data", func(t *testing.T) {
		encryptor, err := New(key)
		if err != nil {
			t.Fatalf("Failed to create encryptor: %v", err)
		}

		invalid := []byte("this is not encrypted data")

		_, err = encryptor.Decrypt(invalid)
		if err == nil {
			t.Error("Expected error for invalid encrypted data, got nil")
		}
	})

	t.Run("Decrypt too short data", func(t *testing.T) {
		encryptor, err := New(key)
		if err != nil {
			t.Fatalf("Failed to create encryptor: %v", err)
		}

		tooShort := []byte("short")

		_, err = encryptor.Decrypt(tooShort)
		if err == nil {
			t.Error("Expected error for too short data, got nil")
		}
	})

	t.Run("Invalid key size", func(t *testing.T) {
		invalidKey := []byte("short")

		_, err := New(invalidKey)
		if err == nil {
			t.Error("Expected error for invalid key size, got nil")
		}
	})

	t.Run("Factory function returns interface", func(t *testing.T) {
		var encryptor encryption.Encryptor
		var err error

		encryptor, err = NewAES(key)
		if err != nil {
			t.Fatalf("NewAES failed: %v", err)
		}

		original := []byte("test data")
		encrypted, err := encryptor.Encrypt(original)
		if err != nil {
			t.Fatalf("Encrypt failed: %v", err)
		}

		decrypted, err := encryptor.Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Decrypt failed: %v", err)
		}

		if !bytes.Equal(original, decrypted) {
			t.Error("Round-trip through interface failed")
		}
	})
}
