package densecode

import (
	"bytes"
	"testing"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func TestEncodeWithoutOptions(t *testing.T) {
	data := []byte("Hello, World! This is a test of densecode without compression or encryption.")

	dc, err := Encode(data, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if dc == nil {
		t.Fatal("Expected non-nil DenseCode")
	}

	if dc.Size < 21 {
		t.Errorf("Expected size >= 21, got %d", dc.Size)
	}
}

func TestEncodeWithCompression(t *testing.T) {
	data := []byte("Hello, World! This is a test of densecode with compression. " +
		"The more data we have, the better compression works. " +
		"Let's add some repetitive text: test test test test test.")

	compressor := &zstd.ZstdCompressor{}

	opts := &Options{
		Compressor: compressor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	dc, err := EncodeWithOptions(data, opts)
	if err != nil {
		t.Fatalf("EncodeWithOptions failed: %v", err)
	}

	if dc == nil {
		t.Fatal("Expected non-nil DenseCode")
	}

	matrix := dc.ToMatrix()
	decoded, err := DecodeWithOptions(matrix, opts)
	if err != nil {
		t.Fatalf("DecodeWithOptions failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original.\nExpected: %s\nGot: %s", data, decoded)
	}
}

func TestEncodeWithEncryption(t *testing.T) {
	data := []byte("Secret message that should be encrypted!")

	key := []byte("12345678901234567890123456789012")
	encryptor, err := aes.NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	opts := &Options{
		Encryptor:  encryptor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	dc, err := EncodeWithOptions(data, opts)
	if err != nil {
		t.Fatalf("EncodeWithOptions failed: %v", err)
	}

	if dc == nil {
		t.Fatal("Expected non-nil DenseCode")
	}

	matrix := dc.ToMatrix()
	decoded, err := DecodeWithOptions(matrix, opts)
	if err != nil {
		t.Fatalf("DecodeWithOptions failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original.\nExpected: %s\nGot: %s", data, decoded)
	}
}

func TestEncodeWithCompressionAndEncryption(t *testing.T) {
	data := []byte("This is a secret message with compression and encryption! " +
		"Let's add more data to make compression worthwhile. " +
		"Repetitive text: test test test test test.")

	compressor := &zstd.ZstdCompressor{}

	key := []byte("12345678901234567890123456789012")
	encryptor, err := aes.NewAES(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	opts := &Options{
		Compressor: compressor,
		Encryptor:  encryptor,
		ErrorLevel: 2,
		ModuleSize: 10,
	}

	dc, err := EncodeWithOptions(data, opts)
	if err != nil {
		t.Fatalf("EncodeWithOptions failed: %v", err)
	}

	if dc == nil {
		t.Fatal("Expected non-nil DenseCode")
	}

	matrix := dc.ToMatrix()
	decoded, err := DecodeWithOptions(matrix, opts)
	if err != nil {
		t.Fatalf("DecodeWithOptions failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original.\nExpected: %s\nGot: %s", data, decoded)
	}
}

func TestEncodeWithWrongDecryptionKey(t *testing.T) {
	data := []byte("Secret message!")

	key1 := []byte("12345678901234567890123456789012")
	encryptor1, err := aes.NewAES(key1)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	opts1 := &Options{
		Encryptor:  encryptor1,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	dc, err := EncodeWithOptions(data, opts1)
	if err != nil {
		t.Fatalf("EncodeWithOptions failed: %v", err)
	}

	key2 := []byte("99999999999999999999999999999999")
	encryptor2, err := aes.NewAES(key2)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	opts2 := &Options{
		Encryptor:  encryptor2,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	matrix := dc.ToMatrix()
	_, err = DecodeWithOptions(matrix, opts2)
	if err == nil {
		t.Error("Expected decryption to fail with wrong key, but it succeeded")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	data := []byte("Test backward compatibility")

	dc, err := Encode(data, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	matrix := dc.ToMatrix()
	decoded, err := Decode(matrix)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original.\nExpected: %s\nGot: %s", data, decoded)
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		errorLevel int
	}{
		{"Short text", []byte("Hello"), 0},
		{"Medium text", []byte("Hello, World! This is a test."), 1},
		{"Long text", []byte("The quick brown fox jumps over the lazy dog. " +
			"Pack my box with five dozen liquor jugs."), 2},
		{"Binary data", []byte{0x00, 0xFF, 0xAA, 0x55, 0x12, 0x34, 0x56, 0x78, 0x90, 0xAB}, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dc, err := Encode(tc.data, tc.errorLevel)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			matrix := dc.ToMatrix()
			decoded, err := Decode(matrix)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if !bytes.Equal(tc.data, decoded) {
				t.Errorf("Decoded data doesn't match original.\nExpected: %v\nGot: %v", tc.data, decoded)
			}
		})
	}
}

func TestNilOptions(t *testing.T) {
	data := []byte("Test with nil options")

	dc, err := EncodeWithOptions(data, nil)
	if err != nil {
		t.Fatalf("EncodeWithOptions with nil failed: %v", err)
	}

	if dc.ModuleSize != 10 {
		t.Errorf("Expected default ModuleSize 10, got %d", dc.ModuleSize)
	}
	if dc.BitsPerModule != 3 {
		t.Errorf("Expected default BitsPerModule 3, got %d", dc.BitsPerModule)
	}

	matrix := dc.ToMatrix()
	decoded, err := DecodeWithOptions(matrix, nil)
	if err != nil {
		t.Fatalf("DecodeWithOptions with nil failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original")
	}
}

func TestCustomModuleSize(t *testing.T) {
	data := []byte("Test custom module size")

	opts := &Options{
		ErrorLevel: 1,
		ModuleSize: 20,
	}

	dc, err := EncodeWithOptions(data, opts)
	if err != nil {
		t.Fatalf("EncodeWithOptions failed: %v", err)
	}

	if dc.ModuleSize != 20 {
		t.Errorf("Expected ModuleSize 20, got %d", dc.ModuleSize)
	}
}

func TestEncodeDecodeWithDifferentBitDensities(t *testing.T) {
	data := []byte("Bit density round-trip test payload")

	for _, bits := range []int{1, 2, 3, 4} {
		t.Run(string(rune('0'+bits)), func(t *testing.T) {
			opts := &Options{
				ErrorLevel:    1,
				ModuleSize:    10,
				BitsPerModule: bits,
			}

			dc, err := EncodeWithOptions(data, opts)
			if err != nil {
				t.Fatalf("EncodeWithOptions failed for bits=%d: %v", bits, err)
			}

			if dc.BitsPerModule != bits {
				t.Fatalf("expected BitsPerModule=%d, got %d", bits, dc.BitsPerModule)
			}

			decoded, err := DecodeWithOptions(dc.ToMatrix(), nil)
			if err != nil {
				t.Fatalf("DecodeWithOptions failed for bits=%d: %v", bits, err)
			}

			if !bytes.Equal(data, decoded) {
				t.Fatalf("decoded data mismatch for bits=%d", bits)
			}
		})
	}
}

func TestInvalidBitsPerModule(t *testing.T) {
	_, err := EncodeWithOptions([]byte("bad"), &Options{
		ErrorLevel:    1,
		ModuleSize:    10,
		BitsPerModule: 5,
	})
	if err == nil {
		t.Fatal("expected error for invalid BitsPerModule, got nil")
	}
}
