package densecode

import (
	"fmt"
)

func (c *Configuration) Decode(matrix [][]int) ([]byte, error) {
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
		decoded, err := c.decodeMatrixWithBits(matrix, bitsPerModule)
		if err == nil {
			return decoded, nil
		}
		lastErr = err
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("unable to decode matrix")
	}
	return nil, lastErr
}

func (c *Configuration) decodeMatrixWithBits(matrix [][]int, bitsPerModule int) ([]byte, error) {
	c.size = len(matrix) // Ensure size is set
	allBytes := c.extractBytesFromMatrix(matrix, bitsPerModule)

	if len(allBytes) < 3 {
		return nil, fmt.Errorf("bits=%d: insufficient data: got %d bytes, need at least 3, matrix size: %dx%d", bitsPerModule, len(allBytes), len(matrix), len(matrix[0]))
	}

	metadata := allBytes[0]
	errorLevel := int((metadata >> 6) & 0x3)
	decodedBitsPerModule := int(((metadata >> 4) & 0x3) + 1)

	if decodedBitsPerModule != bitsPerModule {
		return nil, fmt.Errorf("bits=%d: bits per module mismatch: expected %d, got %d, metadata=%02x", bitsPerModule, bitsPerModule, decodedBitsPerModule, metadata)
	}

	cleanBytes := removeErrorCorrectionBytes(allBytes, errorLevel)

	if len(cleanBytes) < 3 {
		return nil, fmt.Errorf("bits=%d: insufficient data after error correction: %d bytes (had %d, removed %d)", bitsPerModule, len(cleanBytes), len(allBytes), len(allBytes)-len(cleanBytes))
	}

	dataLen := int(cleanBytes[1])<<8 | int(cleanBytes[2])
	if dataLen < 8 || dataLen > 0xFFFF {
		return nil, fmt.Errorf("bits=%d: invalid data length: %d (raw bytes: %02x %02x)", bitsPerModule, dataLen, cleanBytes[1], cleanBytes[2])
	}

	if len(cleanBytes) < 3+dataLen {
		return nil, fmt.Errorf("bits=%d: insufficient data: got %d bytes, need %d (cleanBytes[0]=%02x)", bitsPerModule, len(cleanBytes), 3+dataLen, cleanBytes[0])
	}

	payload := cleanBytes[3 : 3+dataLen]
	processLen := dataLen - 8
	if processLen < 0 {
		return nil, fmt.Errorf("bits=%d: invalid dataLen %d (too small)", bitsPerModule, dataLen)
	}
	processData := payload[:processLen]
	storedChecksum := payload[processLen:]

	decoded, err := c.decryptAndDecompressPayload(processData)
	if err != nil {
		return nil, err
	}

	dataChecksum := c.calculateDataChecksum(decoded)
	if !c.checksumsMatch(storedChecksum, dataChecksum) {
		return nil, fmt.Errorf("bits=%d: checksum mismatch", bitsPerModule)
	}

	return decoded, nil
}

func (c *Configuration) extractBytesFromMatrix(matrix [][]int, bitsPerModule int) []byte {
	size := len(matrix)
	bitBuffer := []int{}

	for layer := 0; layer < (size+1)/2; layer++ {
		for x := layer; x < size-layer; x++ {
			if !c.isReserved(x, layer) {
				value := matrix[x][layer]
				for i := bitsPerModule - 1; i >= 0; i-- {
					bitBuffer = append(bitBuffer, (value>>i)&1)
				}
			}
		}
		for y := layer + 1; y < size-layer; y++ {
			if !c.isReserved(size-layer-1, y) {
				value := matrix[size-layer-1][y]
				for i := bitsPerModule - 1; i >= 0; i-- {
					bitBuffer = append(bitBuffer, (value>>i)&1)
				}
			}
		}
		for x := size - layer - 2; x >= layer; x-- {
			if !c.isReserved(x, size-layer-1) {
				value := matrix[x][size-layer-1]
				for i := bitsPerModule - 1; i >= 0; i-- {
					bitBuffer = append(bitBuffer, (value>>i)&1)
				}
			}
		}
		for y := size - layer - 2; y > layer; y-- {
			if !c.isReserved(layer, y) {
				value := matrix[layer][y]
				for i := bitsPerModule - 1; i >= 0; i-- {
					bitBuffer = append(bitBuffer, (value>>i)&1)
				}
			}
		}
	}

	bytes := make([]byte, (len(bitBuffer)+7)/8)
	for i := 0; i < len(bitBuffer); i++ {
		if i/8 < len(bytes) {
			bytes[i/8] |= byte(bitBuffer[i]) << (7 - (i % 8))
		}
	}

	return bytes
}

func (c *Configuration) decryptAndDecompressPayload(data []byte) ([]byte, error) {
	processed := data

	if c.Encryptor != nil {
		decrypted, err := c.Encryptor.Decrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
		processed = decrypted
	}

	if c.Compressor != nil {
		decompressed, err := c.Compressor.Decompress(processed)
		if err != nil {
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
		processed = decompressed
	}

	return processed, nil
}
