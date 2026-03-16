package densecode

import (
	"crypto/sha256"
	"fmt"
)

type Segment struct {
	Code          *Configuration
	SegmentIndex  int // 0-based index
	TotalSegments int
	dataChecksum  [32]byte // SHA-256 of complete original data
}

type SegmentOptions struct {
	Configuration  *Configuration
	MaxSegmentSize int // Maximum bytes per segment (default: 32KB)
}

func EncodeSegments(data []byte, opts *SegmentOptions) ([]*Segment, error) {
	if opts == nil {
		opts = &SegmentOptions{
			Configuration:  WithDefaults(),
			MaxSegmentSize: 32 * 1024, // 32KB default
		}
	}
	if opts.Configuration == nil {
		opts.Configuration = WithDefaults()
	}
	if opts.MaxSegmentSize == 0 {
		opts.MaxSegmentSize = 32 * 1024
	}

	dataChecksum := sha256.Sum256(data)

	errorLevel := min(max(opts.Configuration.ErrorLevel, 0), 3)
	redundancy := []int{10, 20, 30, 40}[errorLevel]

	effectiveSegmentSize := (opts.MaxSegmentSize*100)/(100+redundancy) - 12
	if effectiveSegmentSize < 1 {
		return nil, fmt.Errorf("segment size too small: %d bytes", effectiveSegmentSize)
	}

	totalSegments := (len(data) + effectiveSegmentSize - 1) / effectiveSegmentSize
	if totalSegments > 65535 {
		return nil, fmt.Errorf("data too large: requires %d segments (max 65535)", totalSegments)
	}

	segments := make([]*Segment, totalSegments)

	for i := range totalSegments {
		start := i * effectiveSegmentSize
		end := min(start+effectiveSegmentSize, len(data))
		segmentdata := data[start:end]

		code, err := encodeSegmentWithOptions(segmentdata, i, totalSegments, dataChecksum, opts.Configuration)
		if err != nil {
			return nil, fmt.Errorf("failed to encode segment %d: %w", i, err)
		}

		segments[i] = &Segment{
			Code:          code,
			SegmentIndex:  i,
			TotalSegments: totalSegments,
			dataChecksum:  dataChecksum,
		}
	}

	return segments, nil
}

func encodeSegmentWithOptions(segmentdata []byte, index, total int, dataChecksum [32]byte, cfg *Configuration) (*Configuration, error) {
	config := WithDefaults()
	config.ModuleSize = cfg.ModuleSize
	config.ErrorLevel = cfg.ErrorLevel
	config.BitsPerModule = cfg.BitsPerModule
	config.Compressor = cfg.Compressor
	config.Encryptor = cfg.Encryptor

	if config.ModuleSize == 0 {
		config.ModuleSize = 10
	}
	bitsPerModule, err := config.normalizeBitsPerModule()
	if err != nil {
		return nil, err
	}

	processed := segmentdata

	if config.Compressor != nil {
		processed, err = config.Compressor.Compress(processed)
		if err != nil {
			return nil, fmt.Errorf("compression failed: %w", err)
		}
	}

	if config.Encryptor != nil {
		processed, err = config.Encryptor.Encrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("encryption failed: %w", err)
		}
	}


	metadata := config.encodeMetadata(bitsPerModule, true)

	dataLen := 1 + len(processed) + 4
	if dataLen > 0xFFFF {
		return nil, fmt.Errorf("segment data too large: %d bytes (max 65535)", dataLen)
	}

	payload := make([]byte, 0, 1+2+2+2+dataLen)
	payload = append(payload, metadata...)
	payload = append(payload, byte(dataLen>>8), byte(dataLen&0xFF))
	payload = append(payload, byte(index>>8), byte(index&0xFF))
	payload = append(payload, byte(total>>8), byte(total&0xFF))
	payload = append(payload, dataChecksum[0]) // First byte of data checksum for verification
	payload = append(payload, processed...)

	segmentHash := sha256.Sum256(segmentdata)
	payload = append(payload, segmentHash[:4]...)

	encoded := config.addErrorCorrection(payload, config.ErrorLevel)

	size, err := config.calculateGridSize(len(encoded), bitsPerModule)
	if err != nil {
		return nil, err
	}

	segmentConfig := &Configuration{
		data:          encoded,
		size:          size,
		ModuleSize:    config.ModuleSize,
		ErrorLevel:    config.ErrorLevel,
		BitsPerModule: bitsPerModule,
		Compressor:    config.Compressor,
		Encryptor:     config.Encryptor,
	}

	return segmentConfig, nil
}

func DecodeSegments(segments []*Segment, cfg *Configuration) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments provided")
	}

	totalSegments := segments[0].TotalSegments
	dataChecksum := segments[0].dataChecksum

	for i, seg := range segments {
		if seg.TotalSegments != totalSegments {
			return nil, fmt.Errorf("segment %d has mismatched total count: expected %d, got %d",
				i, totalSegments, seg.TotalSegments)
		}
		if seg.dataChecksum != dataChecksum {
			return nil, fmt.Errorf("segment %d has mismatched data checksum", i)
		}
	}

	if len(segments) != totalSegments {
		return nil, fmt.Errorf("incomplete segments: have %d, need %d", len(segments), totalSegments)
	}

	segmentMap := make(map[int][]byte)
	config := WithDefaults()
	if cfg != nil {
		config.Compressor = cfg.Compressor
		config.Encryptor = cfg.Encryptor
		config.ErrorLevel = cfg.ErrorLevel
		config.ModuleSize = cfg.ModuleSize
		config.BitsPerModule = cfg.BitsPerModule
	}

	for _, seg := range segments {
		if seg.SegmentIndex < 0 || seg.SegmentIndex >= totalSegments {
			return nil, fmt.Errorf("invalid segment index: %d (total: %d)", seg.SegmentIndex, totalSegments)
		}
		if _, exists := segmentMap[seg.SegmentIndex]; exists {
			return nil, fmt.Errorf("duplicate segment index: %d", seg.SegmentIndex)
		}

		matrix := seg.Code.ToMatrix()
		segmentdata, err := config.decodeSegment(matrix)
		if err != nil {
			return nil, fmt.Errorf("failed to decode segment %d: %w", seg.SegmentIndex, err)
		}

		segmentMap[seg.SegmentIndex] = segmentdata
	}

	var result []byte
	for i := range totalSegments {
		segmentdata, ok := segmentMap[i]
		if !ok {
			return nil, fmt.Errorf("missing segment %d", i)
		}
		result = append(result, segmentdata...)
	}

	resultChecksum := sha256.Sum256(result)
	if resultChecksum != dataChecksum {
		return nil, fmt.Errorf("data checksum mismatch after reassembly")
	}

	return result, nil
}

func (c *Configuration) decodeSegment(matrix [][]int) ([]byte, error) {
	bitsToTry := []int{4, 3, 2, 1}
	if c.BitsPerModule != 0 {
		bitsPerModule, err := c.normalizeBitsPerModule()
		if err != nil {
			return nil, err
		}
		bitsToTry = []int{bitsPerModule}
	}

	c.size = len(matrix)

	var lastErr error
	for _, bitsPerModule := range bitsToTry {
		decoded, err := c.decodeSegmentWithBits(matrix, bitsPerModule)
		if err == nil {
			return decoded, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("unable to decode segment")
	}
	return nil, lastErr
}

func (c *Configuration) decodeSegmentWithBits(matrix [][]int, bitsPerModule int) ([]byte, error) {
	c.size = len(matrix) // Ensure size is set
	bits := c.extractBytesFromMatrix(matrix, bitsPerModule)
	if len(bits) < 8 {
		return nil, fmt.Errorf("insufficient data: got %d bytes, need at least 8", len(bits))
	}

	metadata := bits[0]
	errorLevel, storedBits, isSegment, metadataErr := c.decodeMetadata(metadata, bitsPerModule)
	if metadataErr != nil {
		return nil, metadataErr
	}
	if storedBits != bitsPerModule {
		return nil, fmt.Errorf("bit density mismatch: metadata=%d extracted=%d", storedBits, bitsPerModule)
	}

	if !isSegment {
		return nil, fmt.Errorf("not a segment (use Decode for non-segmented data)")
	}

	dataLen := int(bits[1])<<8 | int(bits[2])
	redundancy := []int{10, 20, 30, 40}[errorLevel]
	totalNeeded := 7 + dataLen

	encodedLen := min((totalNeeded*(100+redundancy)+99)/100, len(bits))
	bits = bits[:encodedLen]

	originalLen := min(max((len(bits)*100)/(100+redundancy), totalNeeded), len(bits))
	bits = bits[:originalLen]

	if len(bits) < totalNeeded {
		return nil, fmt.Errorf("insufficient data: need %d bytes, have %d", totalNeeded, len(bits))
	}

	bits = bits[7:totalNeeded]
	if len(bits) < 1 {
		return nil, fmt.Errorf("insufficient data for data checksum byte")
	}
	bits = bits[1:]
	if len(bits) < 4 {
		return nil, fmt.Errorf("insufficient data for checksum")
	}

	checksumStart := len(bits) - 4
	checksum := bits[checksumStart:]
	processed := bits[:checksumStart]

	var err error
	if c.Encryptor != nil {
		processed, err = c.Encryptor.Decrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
	}
	if c.Compressor != nil {
		processed, err = c.Compressor.Decompress(processed)
		if err != nil {
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
	}

	hash := sha256.Sum256(processed)
	for i := range 4 {
		if checksum[i] != hash[i] {
			return nil, fmt.Errorf("segment checksum mismatch at byte %d: expected %02x, got %02x", i, hash[i], checksum[i])
		}
	}

	return processed, nil
}

func (c *Configuration) decodeMetadata(metadata byte, expectedBits int) (int, int, bool, error) {
	errorLevel := int((metadata >> 6) & 0x3)
	storedBits := int(((metadata >> 4) & 0x3) + 1)
	isSegment := (metadata & 0x08) != 0

	return errorLevel, storedBits, isSegment, nil
}

func DecodeSegmentMetadata(matrix [][]int) (index, total int, dataChecksum byte, err error) {
	config := &Configuration{size: len(matrix)}
	for _, bitsPerModule := range []int{3, 2, 1} {
		bits := config.extractBytesFromMatrix(matrix, bitsPerModule)
		if len(bits) < 8 {
			continue
		}

		metadata := bits[0]
		_, storedBits, isSegment, metadataErr := config.decodeMetadata(metadata, bitsPerModule)
		if metadataErr != nil {
			continue
		}
		if storedBits != bitsPerModule {
			continue
		}

		if !isSegment {
			continue
		}

		index = int(bits[3])<<8 | int(bits[4])
		total = int(bits[5])<<8 | int(bits[6])
		dataChecksum = bits[7]
		return index, total, dataChecksum, nil
	}

	return 0, 0, 0, fmt.Errorf("insufficient or invalid data for segment metadata")
}
