package zstd

import (
	"bytes"
	"context"
	"testing"
	"time"
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

func TestZstdContextCompressor(t *testing.T) {
	compressor := NewContextCompressor()

	t.Run("Compress and Decompress with context", func(t *testing.T) {
		original := []byte("This is test data for context-aware compression")

		ctx := context.Background()

		compressed, err := compressor.CompressWithContext(ctx, original)
		if err != nil {
			t.Fatalf("CompressWithContext failed: %v", err)
		}

		decompressed, err := compressor.DecompressWithContext(ctx, compressed)
		if err != nil {
			t.Fatalf("DecompressWithContext failed: %v", err)
		}

		if !bytes.Equal(original, decompressed) {
			t.Errorf("Decompressed data doesn't match original")
		}
	})

	t.Run("Compress with cancelled context", func(t *testing.T) {
		original := []byte("test data")

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := compressor.CompressWithContext(ctx, original)
		if err == nil {
			t.Error("Expected context cancellation error, got nil")
		}
	})

	t.Run("Decompress with cancelled context", func(t *testing.T) {
		original := []byte("test data")
		compressed, _ := compressor.Compress(original)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := compressor.DecompressWithContext(ctx, compressed)
		if err == nil {
			t.Error("Expected context cancellation error, got nil")
		}
	})

	t.Run("Compress with timeout context", func(t *testing.T) {
		original := []byte("test data for timeout")

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		compressed, err := compressor.CompressWithContext(ctx, original)
		if err != nil {
			t.Fatalf("CompressWithContext with timeout failed: %v", err)
		}

		decompressed, err := compressor.DecompressWithContext(ctx, compressed)
		if err != nil {
			t.Fatalf("DecompressWithContext with timeout failed: %v", err)
		}

		if !bytes.Equal(original, decompressed) {
			t.Error("Decompressed data doesn't match original")
		}
	})
}

func TestContextCompressorInterface(t *testing.T) {
	var _ ContextCompressor = &ZstdContextCompressor{}
}
