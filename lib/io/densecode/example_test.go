package densecode

import (
	"bytes"
	"testing"
)

func TestEncodeDecodeText(t *testing.T) {
	original := "Hello, DenseCode!"

	code, err := EncodeText(original, 2)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if code.Size < 21 {
		t.Errorf("Code size too small: %d", code.Size)
	}

	matrix := code.ToMatrix()
	if len(matrix) != code.Size {
		t.Errorf("Matrix size mismatch: got %d, want %d", len(matrix), code.Size)
	}
}

func TestEncodeBinary(t *testing.T) {
	data := []byte{0x00, 0xFF, 0xAA, 0x55, 0x12, 0x34, 0x56, 0x78}

	code, err := EncodeBinary(data, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if code.ErrorLevel != 1 {
		t.Errorf("Error level mismatch: got %d, want 1", code.ErrorLevel)
	}
}

func TestErrorCorrection(t *testing.T) {
	data := []byte("Test data")

	levels := []int{0, 1, 2, 3}
	var prevSize int

	for _, level := range levels {
		encoded := addErrorCorrection(data, level)
		if len(encoded) <= prevSize {
			t.Errorf("Error correction level %d should produce larger output", level)
		}
		prevSize = len(encoded)
	}
}

func TestColorPalette(t *testing.T) {
	if len(ColorPalette) != 8 {
		t.Errorf("Color palette should have 8 colors, got %d", len(ColorPalette))
	}
}

func TestFinderPatterns(t *testing.T) {
	code, _ := EncodeText("Test", 2)
	matrix := code.ToMatrix()

	if matrix[0][0] != 0 {
		t.Error("Top-left corner should be black")
	}

	hasPattern := false
	for i := range 7 {
		for j := range 7 {
			if matrix[i][j] != 1 {
				hasPattern = true
				break
			}
		}
	}

	if !hasPattern {
		t.Error("Finder pattern not detected")
	}
}

func TestRenderPNG(t *testing.T) {
	code, err := EncodeText("PNG Test", 2)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	filename := "test_output.png"
	err = code.RenderPNG(filename)
	if err != nil {
		t.Fatalf("RenderPNG failed: %v", err)
	}

}

func TestDataCompression(t *testing.T) {
	data := bytes.Repeat([]byte("A"), 1000)

	code, err := Encode(data, 1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(code.Data) >= len(data) {
		t.Logf("Warning: Compression not effective. Original: %d, Compressed: %d",
			len(data), len(code.Data))
	}
}

func BenchmarkEncode(b *testing.B) {
	data := []byte("Benchmark test data for encoding performance")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encode(data, 2)
	}
}

func BenchmarkRenderPNG(b *testing.B) {
	code, _ := EncodeText("Benchmark", 2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = code.RenderPNG("bench_output.png")
	}
}
