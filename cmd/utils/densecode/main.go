package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/io/densecode"
	aesencryption "github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}

	data, err := loadInput(cfg.inputText, cfg.inputFile)
	if err != nil {
		log.Fatalf("Failed to load input: %v", err)
	}

	opts, err := buildOptions(cfg)
	if err != nil {
		log.Fatalf("Failed to build options: %v", err)
	}

	segments, err := densecode.EncodeSegments(data, &densecode.SegmentOptions{
		Options:        opts,
		MaxSegmentSize: cfg.maxSegmentSize,
	})
	if err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}

	fmt.Printf("Input bytes: %d\n", len(data))
	fmt.Printf("Segments: %d\n", len(segments))
	fmt.Printf("Bits/module: %d\n", opts.BitsPerModule)
	fmt.Printf("Error level: %d\n", opts.ErrorLevel)
	fmt.Printf("Compression: %t\n", cfg.enableCompression)
	fmt.Printf("Encryption: %t\n", cfg.encryptionKey != "")

	if cfg.writePNG {
		if err := writeSegmentsPNG(segments, cfg.outputPath); err != nil {
			log.Fatalf("Failed to write PNG: %v", err)
		}
	}

	if cfg.terminal {
		if err := displaySegmentsTerminal(segments, cfg.advanceKey); err != nil {
			log.Fatalf("Terminal display failed: %v", err)
		}
	}
}

type config struct {
	inputText         string
	inputFile         string
	outputPath        string
	terminal          bool
	writePNG          bool
	errorLevel        int
	moduleSize        int
	bits              int
	profile           string
	maxSegmentSize    int
	enableCompression bool
	encryptionKey     string
	advanceKey        string
}

func parseFlags() (*config, error) {
	cfg := &config{}

	flag.StringVar(&cfg.inputText, "text", "", "Text payload to encode")
	flag.StringVar(&cfg.inputFile, "file", "", "File payload to encode")
	flag.StringVar(&cfg.outputPath, "output", "densecode.png", "Output PNG path (single segment) or base name (multi-segment)")
	flag.BoolVar(&cfg.terminal, "terminal", false, "Display code(s) in terminal")
	flag.BoolVar(&cfg.writePNG, "png", true, "Write PNG output")
	flag.IntVar(&cfg.errorLevel, "error", 3, "Error correction level (0-3)")
	flag.IntVar(&cfg.moduleSize, "module-size", 10, "Module size in pixels")
	flag.IntVar(&cfg.bits, "bits", 0, "Bits per module override (1-4). 0 = profile default")
	flag.StringVar(&cfg.profile, "profile", "camera-safe", "Density profile: camera-safe|balanced|max")
	flag.IntVar(&cfg.maxSegmentSize, "max-segment-bytes", 32*1024, "Target max bytes per encoded segment")
	flag.BoolVar(&cfg.enableCompression, "compress", false, "Enable zstd compression")
	flag.StringVar(&cfg.encryptionKey, "encrypt-key-file", "", "Path to file containing AES key (16, 24, or 32 bytes) - prevents exposure in process listings")
	flag.StringVar(&cfg.advanceKey, "advance-key", "enter", "Terminal segment advance key: enter|space|any or custom token")

	flag.Parse()

	if (cfg.inputText != "") == (cfg.inputFile != "") {
		return nil, errors.New("provide exactly one of -text or -file")
	}
	if !cfg.writePNG && !cfg.terminal {
		return nil, errors.New("at least one output mode must be enabled (-png or -terminal)")
	}
	if cfg.errorLevel < 0 || cfg.errorLevel > 3 {
		return nil, errors.New("-error must be in range 0..3")
	}
	if cfg.moduleSize < 1 {
		return nil, errors.New("-module-size must be >= 1")
	}
	if cfg.maxSegmentSize < 256 {
		return nil, errors.New("-max-segment-bytes must be >= 256")
	}

	cfg.profile = strings.ToLower(strings.TrimSpace(cfg.profile))
	cfg.advanceKey = strings.ToLower(strings.TrimSpace(cfg.advanceKey))
	if cfg.advanceKey == "" {
		cfg.advanceKey = "enter"
	}

	return cfg, nil
}

func loadInput(text, path string) ([]byte, error) {
	if text != "" {
		return []byte(text), nil
	}
	if err := validateFilePath(path); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func validateFilePath(path string) error {
	if path == "" {
		return errors.New("file path cannot be empty")
	}
	if strings.Contains(path, "..") {
		return errors.New("path traversal detected: path cannot contain .. components")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("file not accessible: %w", err)
	}
	if info.IsDir() {
		return errors.New("path points to a directory, not a file")
	}
	return nil
}

func buildOptions(cfg *config) (*densecode.Options, error) {
	bits := cfg.bits
	if bits == 0 {
		var err error
		bits, err = bitsForProfile(cfg.profile)
		if err != nil {
			return nil, err
		}
	}

	opts := &densecode.Options{
		ErrorLevel:    cfg.errorLevel,
		ModuleSize:    cfg.moduleSize,
		BitsPerModule: bits,
	}

	if cfg.enableCompression {
		opts.Compressor = &zstd.ZstdCompressor{}
	}

	if cfg.encryptionKey != "" {
		encryptor, err := aesencryption.FromFile(cfg.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load encryption key from file: %w", err)
		}
		opts.Encryptor = encryptor
	}

	return opts, nil
}

func bitsForProfile(profile string) (int, error) {
	switch profile {
	case "camera-safe":
		return 2, nil
	case "balanced":
		return 3, nil
	case "max":
		return 4, nil
	default:
		return 0, fmt.Errorf("invalid -profile %q (allowed: camera-safe|balanced|max)", profile)
	}
}

func writeSegmentsPNG(segments []*densecode.Segment, outputPath string) error {
	if len(segments) == 1 {
		if err := segments[0].Code.RenderPNG(outputPath); err != nil {
			return err
		}
		fmt.Printf("Wrote 1 image: %s\n", outputPath)
		return nil
	}

	for i, seg := range segments {
		name := segmentedFilename(outputPath, i, len(segments))
		if err := seg.Code.RenderPNG(name); err != nil {
			return err
		}
		fmt.Printf("Wrote segment %d/%d: %s\n", i+1, len(segments), name)
	}

	return nil
}

func segmentedFilename(outputPath string, index, total int) string {
	ext := filepath.Ext(outputPath)
	base := strings.TrimSuffix(outputPath, ext)
	if ext == "" {
		ext = ".png"
	}
	return fmt.Sprintf("%s-%04d-of-%04d%s", base, index, total, ext)
}

func displaySegmentsTerminal(segments []*densecode.Segment, advanceKey string) error {
	reader := bufio.NewReader(os.Stdin)

	for i, seg := range segments {
		fmt.Printf("\nSegment %d/%d\n", i+1, len(segments))
		seg.Code.RenderTerminal()

		if i == len(segments)-1 {
			continue
		}

		if err := waitForAdvance(reader, advanceKey, i+2, len(segments)); err != nil {
			return err
		}
	}

	return nil
}

func waitForAdvance(reader *bufio.Reader, advanceKey string, nextSegment, total int) error {
	displayKey := advanceKey
	if displayKey == "" || displayKey == "enter" {
		displayKey = "Enter"
	}

	for {
		fmt.Printf("Press %s for segment %d/%d: ", displayKey, nextSegment, total)
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		token := strings.ToLower(strings.TrimRight(line, "\r\n"))

		switch advanceKey {
		case "", "enter":
			if token == "" {
				return nil
			}
		case "any":
			return nil
		case "space":
			if token == "space" || strings.TrimSpace(token) == "" {
				return nil
			}
		default:
			if token == advanceKey {
				return nil
			}
		}

		fmt.Printf("Input %q did not match %q\n", token, advanceKey)
	}
}
