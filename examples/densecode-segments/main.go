package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/densecode"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func main() {
	fmt.Println("=== DenseCode Multi-Segment Examples ===")

	example1()

	example2()

	example3()

	example4()
}

func example1() {
	fmt.Println("Example 1: Basic Segmentation")
	fmt.Println("------------------------------")

	data := bytes.Repeat([]byte("This is test data. "), 5000) // ~100KB

	opts := &densecode.SegmentOptions{
		Options: &densecode.Options{
			ErrorLevel: 1,
			ModuleSize: 10,
		},
		MaxSegmentSize: 32 * 1024, // 32KB segments
	}

	segments, err := densecode.EncodeSegments(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %d bytes\n", len(data))
	fmt.Printf("Segments created: %d\n", len(segments))
	fmt.Printf("Average segment size: %d bytes\n", len(data)/len(segments))

	for i, seg := range segments {
		fmt.Printf("  Segment %d: %dx%d matrix\n", i, seg.Code.Size, seg.Code.Size)
	}

	decoded, err := densecode.DecodeSegments(segments, opts.Options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded: %d bytes\n", len(decoded))
	fmt.Printf("Match: %v\n\n", bytes.Equal(data, decoded))
}

func example2() {
	fmt.Println("Example 2: Large Data with Compression")
	fmt.Println("---------------------------------------")

	data := bytes.Repeat([]byte("Repetitive content for compression testing. "), 11000) // ~500KB

	compressor := &zstd.ZstdCompressor{}
	opts := &densecode.SegmentOptions{
		Options: &densecode.Options{
			Compressor: compressor,
			ErrorLevel: 1,
			ModuleSize: 10,
		},
		MaxSegmentSize: 32 * 1024,
	}

	segments, err := densecode.EncodeSegments(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %d bytes\n", len(data))
	fmt.Printf("Segments created: %d\n", len(segments))
	fmt.Printf("Compression ratio: %.2f:1\n", float64(len(data))/float64(len(segments)*32*1024))

	decoded, err := densecode.DecodeSegments(segments, opts.Options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded successfully: %v\n", len(decoded) == len(data))
	fmt.Printf("Match: %v\n\n", bytes.Equal(data, decoded))
}

func example3() {
	fmt.Println("Example 3: Segments with Encryption")
	fmt.Println("------------------------------------")

	data := bytes.Repeat([]byte("Secret information. "), 10000) // ~200KB

	key := []byte("12345678901234567890123456789012")
	encryptor, err := aes.NewAES(key)
	if err != nil {
		log.Fatal(err)
	}

	opts := &densecode.SegmentOptions{
		Options: &densecode.Options{
			Encryptor:  encryptor,
			ErrorLevel: 2,
			ModuleSize: 10,
		},
		MaxSegmentSize: 32 * 1024,
	}

	segments, err := densecode.EncodeSegments(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original data: %d bytes\n", len(data))
	fmt.Printf("Encrypted segments: %d\n", len(segments))

	decoded, err := densecode.DecodeSegments(segments, opts.Options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded successfully: %v\n", len(decoded) == len(data))
	fmt.Printf("Match: %v\n\n", bytes.Equal(data, decoded))
}

func example4() {
	fmt.Println("Example 4: Extracting Segment Metadata")
	fmt.Println("---------------------------------------")

	data := bytes.Repeat([]byte("Test. "), 5000) // ~30KB

	opts := &densecode.SegmentOptions{
		Options: &densecode.Options{
			ErrorLevel: 1,
			ModuleSize: 10,
		},
		MaxSegmentSize: 10 * 1024, // 10KB segments
	}

	segments, err := densecode.EncodeSegments(data, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created %d segments\n\n", len(segments))

	for i, seg := range segments {
		matrix := seg.Code.ToMatrix()
		index, total, checksum, err := densecode.DecodeSegmentMetadata(matrix)
		if err != nil {
			log.Printf("Failed to extract metadata from segment %d: %v\n", i, err)
			continue
		}

		fmt.Printf("Segment %d metadata:\n", i)
		fmt.Printf("  Index: %d\n", index)
		fmt.Printf("  Total: %d\n", total)
		fmt.Printf("  Data checksum byte: 0x%02x\n", checksum)
		fmt.Printf("  Matrix size: %dx%d\n", seg.Code.Size, seg.Code.Size)
	}

	fmt.Println()
}
