package zstd

import (
	"bytes"
	"testing"
)

func TestZstdCompressor(t *testing.T) {
	compressor := NewCompressor()

	t.Run("Compress and Decompress", func(t *testing.T) {
		original := []byte("This is test data that should be compressed and decompressed successfully")

		compressed, err := compressor.Compress(original)
		if err != nil {
			t.Fatalf("Compress failed: %v", err)
		}

		if len(compressed) >= len(original) {
			t.Logf("Warning: Compressed size (%d) >= original size (%d)", len(compressed), len(original))
		}

		decompressed, err := compressor.Decompress(compressed)
		if err != nil {
			t.Fatalf("Decompress failed: %v", err)
		}

		if !bytes.Equal(original, decompressed) {
			t.Errorf("Decompressed data doesn't match original.\nOriginal: %s\nDecompressed: %s", original, decompressed)
		}
	})

	t.Run("Compress empty data", func(t *testing.T) {
		original := []byte{}

		compressed, err := compressor.Compress(original)
		if err != nil {
			t.Fatalf("Compress empty data failed: %v", err)
		}

		decompressed, err := compressor.Decompress(compressed)
		if err != nil {
			t.Fatalf("Decompress empty data failed: %v", err)
		}

		if !bytes.Equal(original, decompressed) {
			t.Error("Empty data round-trip failed")
		}
	})

	t.Run("Decompress invalid data", func(t *testing.T) {
		invalid := []byte("this is not compressed data")

		_, err := compressor.Decompress(invalid)
		if err == nil {
			t.Error("Expected error for invalid compressed data, got nil")
		}
	})
}
