package densecode

import (
	"bytes"
	"testing"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func TestEncodeDecodeSegments_Small(t *testing.T) {
	data := []byte("Hello, World!")

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match.\nExpected: %s\nGot: %s", data, decoded)
	}
}

func TestEncodeDecodeSegments_Large(t *testing.T) {
	data := bytes.Repeat([]byte("This is a test message. "), 2000) // ~48KB

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 16 * 1024, // 16KB segments
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	t.Logf("Data size: %d bytes, Segments: %d", len(data), len(segments))

	if len(segments) < 2 {
		t.Errorf("Expected multiple segments, got %d", len(segments))
	}

	for i, seg := range segments {
		if seg.SegmentIndex != i {
			t.Errorf("Segment %d has wrong index: %d", i, seg.SegmentIndex)
		}
		if seg.TotalSegments != len(segments) {
			t.Errorf("Segment %d has wrong total: %d (expected %d)", i, seg.TotalSegments, len(segments))
		}
	}

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match (len: expected %d, got %d)", len(data), len(decoded))
	}
}

func TestEncodeDecodeSegments_OutOfOrder(t *testing.T) {
	data := bytes.Repeat([]byte("Test data. "), 1000) // ~11KB

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 4 * 1024, // 4KB segments
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	if len(segments) < 2 {
		t.Skip("Need at least 2 segments for this test")
	}

	shuffled := make([]*Segment, len(segments))
	copy(shuffled, segments)
	for i, j := 0, len(shuffled)-1; i < j; i, j = i+1, j-1 {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	decoded, err := DecodeSegments(shuffled, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match")
	}
}

func TestEncodeDecodeSegments_WithCompression(t *testing.T) {
	data := bytes.Repeat([]byte("Repetitive content. "), 2000) // ~40KB

	compressor := &zstd.ZstdCompressor{}
	cfg := &Configuration{
		Compressor: compressor,
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 16 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	t.Logf("Original: %d bytes, Segments: %d", len(data), len(segments))

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match")
	}
}

func TestEncodeDecodeSegments_WithEncryption(t *testing.T) {
	data := bytes.Repeat([]byte("Secret message. "), 1000) // ~16KB

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
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 8 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match")
	}
}

func TestEncodeDecodeSegments_WithBoth(t *testing.T) {
	data := bytes.Repeat([]byte("Secret repetitive content. "), 1500) // ~40KB

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
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 16 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	t.Logf("Original: %d bytes, Segments: %d", len(data), len(segments))

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match")
	}
}

func TestDecodeSegmentMetadata(t *testing.T) {
	data := bytes.Repeat([]byte("Test. "), 1000)

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 4 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	for i, seg := range segments {
		matrix := seg.Code.ToMatrix()
		index, total, checksum, err := DecodeSegmentMetadata(matrix)
		if err != nil {
			t.Errorf("DecodeSegmentMetadata failed for segment %d: %v", i, err)
			continue
		}

		if index != i {
			t.Errorf("Segment %d: wrong index %d", i, index)
		}
		if total != len(segments) {
			t.Errorf("Segment %d: wrong total %d (expected %d)", i, total, len(segments))
		}
		if checksum != seg.dataChecksum[0] {
			t.Errorf("Segment %d: wrong checksum byte", i)
		}
	}
}

func TestEncodeDecodeSegments_MaxDensity(t *testing.T) {
	data := bytes.Repeat([]byte("High-density segment test payload. "), 800)

	cfg := &Configuration{
		ErrorLevel:    2,
		ModuleSize:    10,
		BitsPerModule: 4,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 8 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match")
	}
}

func TestEncodeSegments_MissingSegment(t *testing.T) {
	data := bytes.Repeat([]byte("Test. "), 1000)

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 4 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	if len(segments) < 2 {
		t.Skip("Need at least 2 segments for this test")
	}

	incomplete := segments[:len(segments)-1]

	_, err = DecodeSegments(incomplete, cfg)
	if err == nil {
		t.Error("Expected error for incomplete segments, got nil")
	}
}

func TestEncodeSegments_DuplicateSegment(t *testing.T) {
	data := bytes.Repeat([]byte("Test. "), 1000)

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 4 * 1024,
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	if len(segments) < 2 {
		t.Skip("Need at least 2 segments for this test")
	}

	duplicate := append(segments, segments[0])

	_, err = DecodeSegments(duplicate, cfg)
	if err == nil {
		t.Error("Expected error for duplicate segment, got nil")
	}
}

func TestEncodeSegments_VeryLarge(t *testing.T) {
	data := bytes.Repeat([]byte("Large data test. "), 60000) // ~1MB

	cfg := &Configuration{
		ErrorLevel: 1,
		ModuleSize: 10,
	}
	opts := &SegmentOptions{
		Configuration:  cfg,
		MaxSegmentSize: 32 * 1024, // 32KB segments
	}

	segments, err := EncodeSegments(data, opts)
	if err != nil {
		t.Fatalf("EncodeSegments failed: %v", err)
	}

	t.Logf("Data size: %d bytes, Segments: %d", len(data), len(segments))

	decoded, err := DecodeSegments(segments, cfg)
	if err != nil {
		t.Fatalf("DecodeSegments failed: %v", err)
	}

	if !bytes.Equal(data, decoded) {
		t.Errorf("Decoded data doesn't match (len: expected %d, got %d)", len(data), len(decoded))
	}
}
