package main

import (
	"fmt"
	"log"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	"github.com/walterjwhite/go-code/lib/net/messaging"
	"github.com/walterjwhite/go-code/lib/security/encryption"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func example1() {
	fmt.Println("=== Example 1: Individual Interfaces ===")

	var encryptor encryption.Encryptor
	var compressor compression.Compressor

	key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	compressor = zstd.NewCompressor()

	original := []byte("This is sensitive data that needs to be compressed and encrypted")

	compressed, err := compressor.Compress(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original size: %d bytes\n", len(original))
	fmt.Printf("Compressed size: %d bytes\n", len(compressed))

	encrypted, err := encryptor.Encrypt(compressed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted size: %d bytes\n", len(encrypted))

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		log.Fatal(err)
	}

	decompressed, err := compressor.Decompress(decrypted)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %s\n\n", string(decompressed))
}

func example2() {
	fmt.Println("=== Example 2: MessageProcessor ===")

	key := []byte("12345678901234567890123456789012")
	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	compressor := zstd.NewCompressor()
	serializer := serialization.NewJSONSerializer()

	processor := messaging.NewMessageProcessor(
		serializer,
		compressor,
		encryptor,
		false, // serialization handled separately
		true,  // enable compression
		true,  // enable encryption
	)

	original := []byte("This is a message to be processed")
	fmt.Printf("Original: %s\n", string(original))

	processed, err := processor.Process(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Processed size: %d bytes\n", len(processed))

	unprocessed, err := processor.Unprocess(processed)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Unprocessed: %s\n\n", string(unprocessed))
}

type SecureMessageService struct {
	encryptor  encryption.Encryptor
	compressor compression.Compressor
}

func NewSecureMessageService(enc encryption.Encryptor, comp compression.Compressor) *SecureMessageService {
	return &SecureMessageService{
		encryptor:  enc,
		compressor: comp,
	}
}

func (s *SecureMessageService) SendMessage(message []byte) ([]byte, error) {
	compressed, err := s.compressor.Compress(message)
	if err != nil {
		return nil, fmt.Errorf("compression failed: %w", err)
	}

	encrypted, err := s.encryptor.Encrypt(compressed)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	return encrypted, nil
}

func (s *SecureMessageService) ReceiveMessage(encrypted []byte) ([]byte, error) {
	decrypted, err := s.encryptor.Decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	decompressed, err := s.compressor.Decompress(decrypted)
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}

	return decompressed, nil
}

func example3() {
	fmt.Println("=== Example 3: Dependency Injection ===")

	key := []byte("12345678901234567890123456789012")
	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	compressor := zstd.NewCompressor()

	service := NewSecureMessageService(encryptor, compressor)

	original := []byte("Message sent through service")
	fmt.Printf("Original: %s\n", string(original))

	encrypted, err := service.SendMessage(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted size: %d bytes\n", len(encrypted))

	received, err := service.ReceiveMessage(encrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s\n\n", string(received))
}

type MockEncryptor struct {
	EncryptFunc func([]byte) ([]byte, error)
	DecryptFunc func([]byte) ([]byte, error)
}

func (m *MockEncryptor) Encrypt(data []byte) ([]byte, error) {
	if m.EncryptFunc != nil {
		return m.EncryptFunc(data)
	}
	return data, nil
}

func (m *MockEncryptor) Decrypt(data []byte) ([]byte, error) {
	if m.DecryptFunc != nil {
		return m.DecryptFunc(data)
	}
	return data, nil
}

func example4() {
	fmt.Println("=== Example 4: Testing with Mocks ===")

	mock := &MockEncryptor{
		EncryptFunc: func(data []byte) ([]byte, error) {
			return append([]byte("ENCRYPTED:"), data...), nil
		},
		DecryptFunc: func(data []byte) ([]byte, error) {
			if len(data) > 10 {
				return data[10:], nil // Remove "ENCRYPTED:" prefix
			}
			return data, nil
		},
	}

	compressor := zstd.NewCompressor()
	service := NewSecureMessageService(mock, compressor)

	original := []byte("Test message")
	fmt.Printf("Original: %s\n", string(original))

	encrypted, err := service.SendMessage(original)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted: %s\n", string(encrypted))

	received, err := service.ReceiveMessage(encrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s\n\n", string(received))
}

func main() {
	example1()
	example2()
	example3()
	example4()

	fmt.Println("All examples completed successfully!")
}
