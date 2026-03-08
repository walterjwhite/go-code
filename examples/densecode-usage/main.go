package main

import (
	"fmt"
	"log"
	"os"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/densecode"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func loadEncryptionKey() ([]byte, error) {
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		fmt.Println("Error: ENCRYPTION_KEY environment variable not set. Please set it before running.")
		return nil, fmt.Errorf("encryption key not configured")
	}
	return []byte(keyStr), nil
}

func main() {
	fmt.Println("=== DenseCode Examples ===")

	example1()

	example2()

	example3()

	example4()
}

func example1() {
	fmt.Println("Example 1: Basic Encoding")
	fmt.Println("--------------------------")

	data := []byte("Hello, DenseCode!")

	code, err := densecode.Encode(data, 1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %s (%d bytes)\n", data, len(data))
	fmt.Printf("Matrix size: %dx%d\n", code.Size, code.Size)

	decoded, err := densecode.Decode(code.ToMatrix())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded data: %s\n\n", decoded)
}

func example2() {
	fmt.Println("Example 2: With Compression")
	fmt.Println("----------------------------")

	data := []byte("This is a longer message with some repetitive content. " +
		"Repetitive content compresses well. Repetitive content compresses well. " +
		"Repetitive content compresses well.")

	opts := &densecode.Options{
		Compressor: &zstd.ZstdCompressor{},
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	code, err := densecode.EncodeWithOptions(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %d bytes\n", len(data))
	fmt.Printf("Encoded data: %d bytes\n", len(code.Data))
	fmt.Printf("Matrix size: %dx%d\n", code.Size, code.Size)

	decoded, err := densecode.DecodeWithOptions(code.ToMatrix(), opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded: %d bytes\n\n", len(decoded))
}

func example3() {
	fmt.Println("Example 3: With Encryption")
	fmt.Println("---------------------------")

	data := []byte("This is a secret message!")

	key, err := loadEncryptionKey()
	if err != nil {
		log.Fatal(err)
	}

	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	opts := &densecode.Options{
		Encryptor:  encryptor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	code, err := densecode.EncodeWithOptions(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %s (%d bytes)\n", data, len(data))
	fmt.Printf("Matrix size: %dx%d\n", code.Size, code.Size)

	decoded, err := densecode.DecodeWithOptions(code.ToMatrix(), opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded data: %s\n\n", decoded)
}

func example4() {
	fmt.Println("Example 4: With Compression and Encryption")
	fmt.Println("-------------------------------------------")

	data := []byte("This is a secret message with repetitive content. " +
		"Repetitive content. Repetitive content. Repetitive content.")

	key, err := loadEncryptionKey()
	if err != nil {
		log.Fatal(err)
	}

	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	opts := &densecode.Options{
		Compressor: &zstd.ZstdCompressor{},
		Encryptor:  encryptor,
		ErrorLevel: 2,
		ModuleSize: 10,
	}

	code, err := densecode.EncodeWithOptions(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %d bytes\n", len(data))
	fmt.Printf("Encoded data: %d bytes\n", len(code.Data))
	fmt.Printf("Matrix size: %dx%d\n", code.Size, code.Size)

	decoded, err := densecode.DecodeWithOptions(code.ToMatrix(), opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded: %d bytes\n\n", len(decoded))
}
