package densecode

import (
	"crypto/sha256"
	"fmt"
	"math"
)

func (c *Configuration) encode(data []byte) error {
	bitsPerModule, err := c.normalizeBitsPerModule()
	if err != nil {
		return err
	}

	processed, err := c.processDataPipeline(data)
	if err != nil {
		return err
	}

	payload, err := c.buildPayload(processed, data, bitsPerModule)
	if err != nil {
		return err
	}

	encoded := c.addErrorCorrection(payload, c.ErrorLevel)
	size, err := c.calculateGridSize(len(encoded), bitsPerModule)
	if err != nil {
		return err
	}

	c.data = encoded
	c.size = size
	return nil
}

func (c *Configuration) normalizeBitsPerModule() (int, error) {
	bits := c.BitsPerModule
	if bits == 0 {
		bits = 3
	}
	if bits < 1 || bits > 4 {
		return 0, fmt.Errorf("invalid BitsPerModule %d: must be 1-4", bits)
	}
	return bits, nil
}

func (c *Configuration) processDataPipeline(data []byte) ([]byte, error) {
	processed := data
	if c.Compressor != nil {
		compressed, err := c.Compressor.Compress(processed)
		if err != nil {
			return nil, fmt.Errorf("compression failed: %w", err)
		}
		processed = compressed
	}
	if c.Encryptor != nil {
		encrypted, err := c.Encryptor.Encrypt(processed)
		if err != nil {
			return nil, fmt.Errorf("encryption failed: %w", err)
		}
		processed = encrypted
	}
	return processed, nil
}

func (c *Configuration) buildPayload(processed, originalData []byte, bitsPerModule int) ([]byte, error) {
	hash := sha256.Sum256(originalData)
	metadata := c.encodeMetadata(bitsPerModule, false)
	dataLen := len(processed) + 8
	if dataLen > 0xFFFF {
		return nil, fmt.Errorf("data too large: %d bytes (max 65535)", dataLen)
	}
	payload := make([]byte, 0, 3+len(processed)+8)
	payload = append(payload, metadata...)
	payload = append(payload, byte(dataLen>>8), byte(dataLen&0xFF))
	payload = append(payload, processed...)
	payload = append(payload, hash[:8]...)
	return payload, nil
}

func (c *Configuration) encodeMetadata(bitsPerModule int, isSegment bool) []byte {
	metadata := byte(c.ErrorLevel&0x3)<<6 | byte(bitsPerModule-1)<<4&0x30
	if isSegment {
		metadata |= 0x08
	}
	return []byte{metadata}
}

func (c *Configuration) calculateGridSize(encodedLen, bitsPerModule int) (int, error) {
	bitsNeeded := encodedLen * 8
	modulesNeeded := (bitsNeeded + bitsPerModule - 1) / bitsPerModule

	size := max(int(math.Sqrt(float64(modulesNeeded)))+1, 21)
	for {
		totalModules := size * size
		reservedModules := 3*7*7 + 2*(size-14)
		if totalModules-reservedModules >= modulesNeeded {
			break
		}
		size++
		if size > 1000 {
			return 0, fmt.Errorf("data too large: requires matrix size > 1000")
		}
	}
	return size, nil
}

func (c *Configuration) addErrorCorrection(data []byte, errorLevel int) []byte {
	redundancy := (errorLevel + 1) * 8
	encoded := make([]byte, len(data)+redundancy)
	copy(encoded, data)

	for i := range redundancy {
		parity := byte(0)
		for j := i; j < len(data); j += redundancy {
			parity ^= data[j]
		}
		encoded[len(data)+i] = parity
	}

	return encoded
}

func (c *Configuration) ToMatrix() [][]int {
	matrix := make([][]int, c.size)
	for i := range matrix {
		matrix[i] = make([]int, c.size)
		for j := range matrix[i] {
			matrix[i][j] = 1 // Default to white
		}
	}

	c.addFinderPatterns(matrix)
	c.addTimingPatterns(matrix)
	c.encodeData(matrix)

	return matrix
}

func (c *Configuration) addFinderPatterns(matrix [][]int) {
	c.drawFinderPattern(matrix, 0, 0)
	c.drawFinderPattern(matrix, c.size-7, 0)
	c.drawFinderPattern(matrix, 0, c.size-7)
}

func (c *Configuration) drawFinderPattern(matrix [][]int, x, y int) {
	pattern := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	for i := range 7 {
		for j := range 7 {
			matrix[x+i][y+j] = pattern[i][j]
		}
	}
}

func (c *Configuration) addTimingPatterns(matrix [][]int) {
	maxColor := (1 << c.bitsPerModule()) - 1
	for i := 8; i < c.size-8; i++ {
		matrix[6][i] = (i % 2) * maxColor
		matrix[i][6] = (i % 2) * maxColor
	}
}

func (c *Configuration) encodeData(matrix [][]int) {
	bitPos := 0
	dataLen := len(c.data) * 8
	bitsPerModule := c.bitsPerModule()

	for layer := 0; layer < (c.size+1)/2; layer++ {
		for x := layer; x < c.size-layer && bitPos < dataLen; x++ {
			if !c.isReserved(x, layer) {
				c.setModule(matrix, x, layer, bitPos, bitsPerModule)
				bitPos += bitsPerModule
			}
		}
		for y := layer + 1; y < c.size-layer && bitPos < dataLen; y++ {
			if !c.isReserved(c.size-layer-1, y) {
				c.setModule(matrix, c.size-layer-1, y, bitPos, bitsPerModule)
				bitPos += bitsPerModule
			}
		}
		for x := c.size - layer - 2; x >= layer && bitPos < dataLen; x-- {
			if !c.isReserved(x, c.size-layer-1) {
				c.setModule(matrix, x, c.size-layer-1, bitPos, bitsPerModule)
				bitPos += bitsPerModule
			}
		}
		for y := c.size - layer - 2; y > layer && bitPos < dataLen; y-- {
			if !c.isReserved(layer, y) {
				c.setModule(matrix, layer, y, bitPos, bitsPerModule)
				bitPos += bitsPerModule
			}
		}
	}
}

func (c *Configuration) setModule(matrix [][]int, x, y, bitPos, bitsPerModule int) {
	if bitPos+bitsPerModule > len(c.data)*8 {
		return
	}

	value := 0
	for i := range bitsPerModule {
		byteIdx := (bitPos + i) / 8
		bitIdx := 7 - ((bitPos + i) % 8)
		if byteIdx < len(c.data) {
			value = (value << 1) | int((c.data[byteIdx]>>bitIdx)&1)
		}
	}

	matrix[x][y] = value
}

func (c *Configuration) isReserved(x, y int) bool {
	if (x < 7 && y < 7) || (x >= c.size-7 && y < 7) || (x < 7 && y >= c.size-7) {
		return true
	}
	if x == 6 || y == 6 {
		return true
	}
	return false
}

func (c *Configuration) bitsPerModule() int {
	bits, err := c.normalizeBitsPerModule()
	if err != nil {
		return 3
	}
	return bits
}
