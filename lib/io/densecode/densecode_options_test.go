package densecode

import (
	"bytes"
	"os"
	"testing"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func TestEncodeWithoutOptions(t *testing.T) {
	data := []byte("Hello, World! This is a test of densecode without compression or encryption.")

	cfg := &Configuration{}
	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result == nil || len(result.Segments) == 0 {
		t.Fatal("Expected non-nil result with segments")
	}

	if result.Segments[0].Code.size < 21 {
		t.Errorf("Expected size >= 21, got %d", result.Segments[0].Code.size)
	}
}

func TestEncodeWithCompression(t *testing.T) {
	data := []byte("Hello, World! This is a test of densecode with compression. " +
		"The more data we have, the better compression works. " +
		"Let's add some repetitive text: test test test test test.")

	compressor := &zstd.ZstdCompressor{}

	cfg := &Configuration{
		Compressor: compressor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result == nil || len(result.Segments) == 0 {
		t.Fatal("Expected non-nil result with segments")
	}

	matrix := result.Segments[0].Code.ToMatrix()
	decoded, err := cfg.Decode(matrix)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
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

	cfg := &Configuration{
		Encryptor:  encryptor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result == nil || len(result.Segments) == 0 {
		t.Fatal("Expected non-nil result with segments")
	}

	matrix := result.Segments[0].Code.ToMatrix()
	decoded, err := cfg.Decode(matrix)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
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

	cfg := &Configuration{
		Compressor: compressor,
		Encryptor:  encryptor,
		ErrorLevel: 2,
		ModuleSize: 10,
	}

	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result == nil || len(result.Segments) == 0 {
		t.Fatal("Expected non-nil result with segments")
	}

	matrix := result.Segments[0].Code.ToMatrix()
	decoded, err := cfg.Decode(matrix)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
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

	cfg1 := &Configuration{
		Encryptor:  encryptor1,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	result, err := cfg1.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	key2 := []byte("99999999999999999999999999999999")
	encryptor2, err := aes.NewAES(key2)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	cfg2 := &Configuration{
		Encryptor:  encryptor2,
		ErrorLevel: 1,
		ModuleSize: 10,
	}

	matrix := result.Segments[0].Code.ToMatrix()
	_, err = cfg2.Decode(matrix)
	if err == nil {
		t.Error("Expected decryption to fail with wrong key, but it succeeded")
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
			cfg := &Configuration{ErrorLevel: tc.errorLevel}
			result, err := cfg.Encode(tc.data)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			matrix := result.Segments[0].Code.ToMatrix()
			decoded, err := cfg.Decode(matrix)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			if !bytes.Equal(tc.data, decoded) {
				t.Errorf("Decoded data doesn't match original.\nExpected: %v\nGot: %v", tc.data, decoded)
			}
		})
	}
}

func TestNilConfiguration(t *testing.T) {
	data := []byte("Test with nil configuration")

	cfg := &Configuration{}
	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result.Segments[0].Code.ModuleSize != 10 {
		t.Errorf("Expected default ModuleSize 10, got %d", result.Segments[0].Code.ModuleSize)
	}

	matrix := result.Segments[0].Code.ToMatrix()
	decoded, err := cfg.Decode(matrix)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match original")
	}
}

func TestCustomModuleSize(t *testing.T) {
	data := []byte("Test custom module size")

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 20,
	}

	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if result.Segments[0].Code.ModuleSize != 20 {
		t.Errorf("Expected ModuleSize 20, got %d", result.Segments[0].Code.ModuleSize)
	}
}

func TestEncodeDecodeWithDifferentBitDensities(t *testing.T) {
	data := []byte("Bit density round-trip test payload")

	for _, bits := range []int{1, 2, 3, 4} {
		t.Run(string(rune('0'+bits)), func(t *testing.T) {
			cfg := &Configuration{
				ErrorLevel:    1,
				ModuleSize:    10,
				BitsPerModule: bits,
			}

			result, err := cfg.Encode(data)
			if err != nil {
				t.Fatalf("Encode failed for bits=%d: %v", bits, err)
			}

			if result.Segments[0].Code.BitsPerModule != bits {
				t.Fatalf("expected BitsPerModule=%d, got %d", bits, result.Segments[0].Code.BitsPerModule)
			}

			decoded, err := cfg.Decode(result.Segments[0].Code.ToMatrix())
			if err != nil {
				t.Fatalf("Decode failed for bits=%d: %v", bits, err)
			}

			if !bytes.Equal(data, decoded) {
				t.Fatalf("decoded data mismatch for bits=%d", bits)
			}
		})
	}
}

func TestInvalidBitsPerModule(t *testing.T) {
	cfg := &Configuration{
		ErrorLevel:    1,
		ModuleSize:    10,
		BitsPerModule: 5,
	}
	_, err := cfg.Encode([]byte("bad"))
	if err == nil {
		t.Fatal("expected error for invalid BitsPerModule, got nil")
	}
}

func TestAutomaticSegmentation(t *testing.T) {
	data := make([]byte, 50*1024) // 50KB
	for i := range data {
		data[i] = byte(i % 256)
	}

	cfg := &Configuration{
		MaxSegmentSize: 20 * 1024, // Force smaller segments for testing
	}

	result, err := cfg.Encode(data)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if !result.IsMultiSegment {
		t.Error("Expected multi-segment result for large data")
	}

	if len(result.Segments) < 2 {
		t.Errorf("Expected multiple segments, got %d", len(result.Segments))
	}

	decoded, err := DecodeSegments(result.Segments, cfg)
	if err != nil {
		decoded, err = cfg.Decode(result.Segments[0].Code.ToMatrix())
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
	}

	if !bytes.Equal(data, decoded) {
		t.Error("Decoded data doesn't match original")
	}
}

func TestEncodeFiles(t *testing.T) {
	files := []string{"test_file1.txt", "test_file2.txt"}
	for i, name := range files {
		content := "Content of file " + string(rune('1'+i))
		err := os.WriteFile(name, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer func() { _ = os.Remove(name) }()
	}

	result, err := EncodeFiles(files, &Configuration{})
	if err != nil {
		t.Fatalf("EncodeFiles failed: %v", err)
	}

	if result == nil || len(result.Segments) == 0 {
		t.Fatal("Expected non-nil result with segments")
	}

	cfg := &Configuration{}
	decoded, err := cfg.DecodeFiles(result)
	if err != nil {
		t.Fatalf("DecodeFiles failed: %v", err)
	}

	if len(decoded) == 0 {
		t.Error("Expected non-empty decoded data")
	}
}
