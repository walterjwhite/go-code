package densecode

import (
	"crypto/sha256"
	"fmt"
)

type Segment struct {
	Code          *DenseCode
	SegmentIndex  int // 0-based index
	TotalSegments int
	DataChecksum  [32]byte // SHA-256 of complete original data
}

type SegmentOptions struct {
	*Options
	MaxSegmentSize int // Maximum bytes per segment (default: 32KB)
}

func EncodeSegments(data []byte, opts *SegmentOptions) ([]*Segment, error) {
	if opts == nil {
		opts = &SegmentOptions{
			Options:        &Options{ErrorLevel: 1, ModuleSize: 10},
			MaxSegmentSize: 32 * 1024, // 32KB default
		}
	}
	if opts.Options == nil {
		opts.Options = &Options{ErrorLevel: 1, ModuleSize: 10}
	}
	if opts.MaxSegmentSize == 0 {
		opts.MaxSegmentSize = 32 * 1024
	}

	dataChecksum := sha256.Sum256(data)

	errorLevel := min(max(opts.ErrorLevel, 0), 3)
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
		segmentData := data[start:end]

		code, err := encodeSegment(segmentData, i, totalSegments, dataChecksum, opts.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to encode segment %d: %w", i, err)
		}

		segments[i] = &Segment{
			Code:          code,
			SegmentIndex:  i,
			TotalSegments: totalSegments,
			DataChecksum:  dataChecksum,
		}
	}

	return segments, nil
}

func encodeSegment(segmentData []byte, index, total int, dataChecksum [32]byte, opts *Options) (*DenseCode, error) {
	if opts == nil {
		opts = &Options{ErrorLevel: 1, ModuleSize: 10}
	}
	if opts.ModuleSize == 0 {
		opts.ModuleSize = 10
	}
	bitsPerModule, err := normalizeBitsPerModule(opts.BitsPerModule)
	if err != nil {
		return nil, err
	}

	processed := segmentData

	if opts.Compressor != nil {
		processed, err = opts.Compressor.Compress(processed)
		if err != nil {
			return nil, fmt.Errorf("compression failed: %w", err)
		}
	}

	if opts.Encryptor != nil {
		processed, err = opts.Encryptor.Encrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("encryption failed: %w", err)
		}
	}


	metadata := encodeMetadata(opts.ErrorLevel, bitsPerModule, true)

	dataLen := 1 + len(processed) + 4
	if dataLen > 0xFFFF {
		return nil, fmt.Errorf("segment data too large: %d bytes (max 65535)", dataLen)
	}

	payload := make([]byte, 0, 1+2+2+2+dataLen)
	payload = append(payload, metadata)
	payload = append(payload, byte(dataLen>>8), byte(dataLen&0xFF))
	payload = append(payload, byte(index>>8), byte(index&0xFF))
	payload = append(payload, byte(total>>8), byte(total&0xFF))
	payload = append(payload, dataChecksum[0]) // First byte of data checksum for verification
	payload = append(payload, processed...)

	segmentHash := sha256.Sum256(segmentData)
	payload = append(payload, segmentHash[:4]...)

	encoded := addErrorCorrection(payload, opts.ErrorLevel)

	bitsNeeded := len(encoded) * 8
	modulesNeeded := (bitsNeeded + bitsPerModule - 1) / bitsPerModule

	size := max(int(sqrt(modulesNeeded))+1, 21)

	for {
		totalModules := size * size
		reservedModules := 3*7*7 + 2*(size-14)
		usableModules := totalModules - reservedModules

		if usableModules >= modulesNeeded {
			break
		}
		size++

		if size > 1000 {
			return nil, fmt.Errorf("segment too large: requires matrix size > 1000")
		}
	}

	return &DenseCode{
		Data:          encoded,
		Size:          size,
		ModuleSize:    opts.ModuleSize,
		ErrorLevel:    opts.ErrorLevel,
		BitsPerModule: bitsPerModule,
	}, nil
}

func DecodeSegments(segments []*Segment, opts *Options) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments provided")
	}

	totalSegments := segments[0].TotalSegments
	dataChecksum := segments[0].DataChecksum

	for i, seg := range segments {
		if seg.TotalSegments != totalSegments {
			return nil, fmt.Errorf("segment %d has mismatched total count: expected %d, got %d",
				i, totalSegments, seg.TotalSegments)
		}
		if seg.DataChecksum != dataChecksum {
			return nil, fmt.Errorf("segment %d has mismatched data checksum", i)
		}
	}

	if len(segments) != totalSegments {
		return nil, fmt.Errorf("incomplete segments: have %d, need %d", len(segments), totalSegments)
	}

	segmentMap := make(map[int][]byte)
	for _, seg := range segments {
		if seg.SegmentIndex < 0 || seg.SegmentIndex >= totalSegments {
			return nil, fmt.Errorf("invalid segment index: %d (total: %d)", seg.SegmentIndex, totalSegments)
		}
		if _, exists := segmentMap[seg.SegmentIndex]; exists {
			return nil, fmt.Errorf("duplicate segment index: %d", seg.SegmentIndex)
		}

		matrix := seg.Code.ToMatrix()
		segmentData, err := decodeSegment(matrix, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to decode segment %d: %w", seg.SegmentIndex, err)
		}

		segmentMap[seg.SegmentIndex] = segmentData
	}

	var result []byte
	for i := range totalSegments {
		segmentData, ok := segmentMap[i]
		if !ok {
			return nil, fmt.Errorf("missing segment %d", i)
		}
		result = append(result, segmentData...)
	}

	resultChecksum := sha256.Sum256(result)
	if resultChecksum != dataChecksum {
		return nil, fmt.Errorf("data checksum mismatch after reassembly")
	}

	return result, nil
}

func decodeSegment(matrix [][]int, opts *Options) ([]byte, error) {
	if opts == nil {
		opts = &Options{}
	}
	bitsToTry := []int{4, 3, 2, 1}
	if opts.BitsPerModule != 0 {
		bitsPerModule, err := normalizeBitsPerModule(opts.BitsPerModule)
		if err != nil {
			return nil, err
		}
		bitsToTry = []int{bitsPerModule}
	}

	var lastErr error
	for _, bitsPerModule := range bitsToTry {
		decoded, err := decodeSegmentWithBits(matrix, opts, bitsPerModule)
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

func DecodeSegmentMetadata(matrix [][]int) (index, total int, dataChecksum byte, err error) {
	for _, bitsPerModule := range []int{3, 2, 1} {
		bits := extractBytesFromMatrix(matrix, bitsPerModule)
		if len(bits) < 8 {
			continue
		}

		metadata := bits[0]
		_, storedBits, isSegment, err := decodeMetadata(metadata)
		if err != nil {
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

func decodeSegmentWithBits(matrix [][]int, opts *Options, bitsPerModule int) ([]byte, error) {
	bits := extractBytesFromMatrix(matrix, bitsPerModule)
	if len(bits) < 8 {
		return nil, fmt.Errorf("insufficient data: got %d bytes, need at least 8", len(bits))
	}

	metadata := bits[0]
	errorLevel, storedBits, isSegment, metadataErr := decodeMetadata(metadata)
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
	if opts.Encryptor != nil {
		processed, err = opts.Encryptor.Decrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
	}
	if opts.Compressor != nil {
		processed, err = opts.Compressor.Decompress(processed)
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
