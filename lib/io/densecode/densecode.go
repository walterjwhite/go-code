package densecode

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

type DenseCode struct {
	Data          []byte
	Size          int
	ModuleSize    int
	ErrorLevel    int // 0-3, higher = more redundancy
	BitsPerModule int // 1-4, higher = more density and lower camera tolerance
}

type Options struct {
	Compressor    compression.Compressor // Optional: nil means no compression
	Encryptor     encryption.Encryptor   // Optional: nil means no encryption
	ErrorLevel    int                    // 0-3, higher = more redundancy
	ModuleSize    int                    // Size of each module in pixels (default: 10)
	BitsPerModule int                    // 1-4 bits per module (default: 3)
}

var ColorPalette = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 255, 0, 255},     // Green
	{0, 0, 255, 255},     // Blue
	{255, 255, 0, 255},   // Yellow
	{255, 0, 255, 255},   // Magenta
	{0, 255, 255, 255},   // Cyan
}

var colorPalette1Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
}

var colorPalette2Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 0, 255, 255},     // Blue
}

var colorPalette4Bit = []color.RGBA{
	{0, 0, 0, 255},       // Black
	{255, 255, 255, 255}, // White
	{255, 0, 0, 255},     // Red
	{0, 255, 0, 255},     // Lime
	{0, 0, 255, 255},     // Blue
	{255, 255, 0, 255},   // Yellow
	{255, 0, 255, 255},   // Magenta
	{0, 255, 255, 255},   // Cyan
	{128, 0, 0, 255},     // Maroon
	{0, 128, 0, 255},     // Green
	{0, 0, 128, 255},     // Navy
	{128, 128, 0, 255},   // Olive
	{128, 0, 128, 255},   // Purple
	{0, 128, 128, 255},   // Teal
	{192, 192, 192, 255}, // Silver
	{128, 128, 128, 255}, // Gray
}

func Encode(data []byte, errorLevel int) (*DenseCode, error) {
	return EncodeWithOptions(data, &Options{
		ErrorLevel:    errorLevel,
		ModuleSize:    10,
		BitsPerModule: 3,
	})
}

func EncodeWithOptions(data []byte, opts *Options) (*DenseCode, error) {
	if opts == nil {
		opts = &Options{ErrorLevel: 0, ModuleSize: 10, BitsPerModule: 3}
	}
	if opts.ModuleSize == 0 {
		opts.ModuleSize = 10
	}
	bitsPerModule, bitsErr := normalizeBitsPerModule(opts.BitsPerModule)
	if bitsErr != nil {
		return nil, bitsErr
	}

	processed := data
	var err error

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

	hash := sha256.Sum256(data)

	metadata := encodeMetadata(opts.ErrorLevel, bitsPerModule, false)
	dataLen := len(processed) + 8 // processed data + 8 byte checksum (increased from 4 to 8 for better security)
	if dataLen > 0xFFFF {
		return nil, fmt.Errorf("data too large: %d bytes (max 65535)", dataLen)
	}

	payload := make([]byte, 0, 3+len(processed)+8)
	payload = append(payload, metadata)
	payload = append(payload, byte(dataLen>>8), byte(dataLen&0xFF)) // 2-byte length
	payload = append(payload, processed...)
	payload = append(payload, hash[:8]...) // Use 8 bytes instead of 4 for better integrity verification

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
			return nil, fmt.Errorf("data too large: requires matrix size > 1000")
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

func addErrorCorrection(data []byte, level int) []byte {
	redundancy := []int{10, 20, 30, 40}[level] // percentage
	extraBytes := (len(data) * redundancy) / 100

	result := make([]byte, len(data)+extraBytes)
	copy(result, data)

	for i := range extraBytes {
		parity := byte(0)
		for j := range data {
			parity = rotateLeft(parity, 1) ^ data[j]
			if j%extraBytes == i {
				parity = rotateLeft(parity, 3)
			}
		}
		result[len(data)+i] = parity
	}

	return result
}

func rotateLeft(b byte, n uint) byte {
	n = n % 8
	return (b << n) | (b >> (8 - n))
}

func sqrt(n int) int {
	if n < 2 {
		return n
	}
	x := n
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + n/x) / 2
	}
	return x
}

func (dc *DenseCode) ToMatrix() [][]int {
	matrix := make([][]int, dc.Size)
	for i := range matrix {
		matrix[i] = make([]int, dc.Size)
	}

	dc.addFinderPatterns(matrix)

	dc.addTimingPatterns(matrix)

	dc.encodeData(matrix)

	return matrix
}

func (dc *DenseCode) addFinderPatterns(matrix [][]int) {
	dc.drawFinderPattern(matrix, 0, 0)
	dc.drawFinderPattern(matrix, dc.Size-7, 0)
	dc.drawFinderPattern(matrix, 0, dc.Size-7)
}

func (dc *DenseCode) drawFinderPattern(matrix [][]int, x, y int) {
	pattern := [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 2, 0, 1, 0},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	for i := range 7 {
		for j := range 7 {
			if x+i < dc.Size && y+j < dc.Size {
				matrix[y+j][x+i] = pattern[j][i]
			}
		}
	}
}

func (dc *DenseCode) addTimingPatterns(matrix [][]int) {
	maxColor := (1 << dc.bitsPerModule()) - 1
	for i := 8; i < dc.Size-8; i++ {
		matrix[6][i] = (i % 2) * maxColor
		matrix[i][6] = (i % 2) * maxColor
	}
}

func (dc *DenseCode) encodeData(matrix [][]int) {
	bitPos := 0
	dataLen := len(dc.Data) * 8
	bitsPerModule := dc.bitsPerModule()

	for y := dc.Size - 1; y >= 0; y-- {
		for x := dc.Size - 1; x >= 0; x-- {
			if dc.isReserved(x, y) {
				continue
			}

			if bitPos < dataLen {
				matrix[y][x] = extractPackedBits(dc.Data, bitPos, bitsPerModule)
				bitPos += bitsPerModule
			}
		}
	}
}

func (dc *DenseCode) isReserved(x, y int) bool {
	if (x < 7 && y < 7) || (x >= dc.Size-7 && y < 7) || (x < 7 && y >= dc.Size-7) {
		return true
	}
	if x == 6 || y == 6 {
		return true
	}
	return false
}

func (dc *DenseCode) RenderPNG(filename string) error {
	matrix := dc.ToMatrix()
	imgSize := dc.Size * dc.ModuleSize
	palette := colorPaletteForBits(dc.bitsPerModule())

	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	for y := range imgSize {
		for x := range imgSize {
			img.Set(x, y, color.White)
		}
	}

	for y := 0; y < dc.Size; y++ {
		for x := 0; x < dc.Size; x++ {
			colorIdx := matrix[y][x]
			if colorIdx >= len(palette) {
				colorIdx = 0
			}
			c := palette[colorIdx]

			for dy := 0; dy < dc.ModuleSize; dy++ {
				for dx := 0; dx < dc.ModuleSize; dx++ {
					img.Set(x*dc.ModuleSize+dx, y*dc.ModuleSize+dy, c)
				}
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := os.Chmod(filename, 0600); err != nil {
		_ = f.Close()
		return err
	}
	defer func() { _ = f.Close() }()

	return png.Encode(f, img)
}

func (dc *DenseCode) RenderTerminal() {
	matrix := dc.ToMatrix()
	palette := colorPaletteForBits(dc.bitsPerModule())
	reset := "\033[0m"

	fmt.Println("\nDenseCode:")
	for y := 0; y < dc.Size; y++ {
		for x := 0; x < dc.Size; x++ {
			colorIdx := matrix[y][x]
			if colorIdx >= len(palette) {
				colorIdx = 0
			}
			c := palette[colorIdx]
			fmt.Printf("\033[48;2;%d;%d;%dm  %s", c.R, c.G, c.B, reset)
		}
		fmt.Println()
	}
	fmt.Println()
}

func EncodeText(text string, errorLevel int) (*DenseCode, error) {
	return Encode([]byte(text), errorLevel)
}

func EncodeBinary(data []byte, errorLevel int) (*DenseCode, error) {
	return Encode(data, errorLevel)
}

func EncodeBase64(data []byte, errorLevel int) (*DenseCode, error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	return Encode([]byte(encoded), errorLevel)
}

func Decode(matrix [][]int) ([]byte, error) {
	return DecodeWithOptions(matrix, &Options{})
}

func DecodeWithOptions(matrix [][]int, opts *Options) ([]byte, error) {
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
		decoded, err := decodeMatrixWithBits(matrix, opts, bitsPerModule)
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

func isReservedPos(x, y, size int) bool {
	if (x < 7 && y < 7) || (x >= size-7 && y < 7) || (x < 7 && y >= size-7) {
		return true
	}
	if x == 6 || y == 6 {
		return true
	}
	return false
}

func normalizeBitsPerModule(bits int) (int, error) {
	if bits == 0 {
		return 3, nil
	}
	if bits < 1 || bits > 4 {
		return 0, fmt.Errorf("invalid BitsPerModule %d: must be 1-4", bits)
	}
	return bits, nil
}

func (dc *DenseCode) bitsPerModule() int {
	bits, err := normalizeBitsPerModule(dc.BitsPerModule)
	if err != nil {
		return 3
	}
	return bits
}

func colorPaletteForBits(bitsPerModule int) []color.RGBA {
	switch bitsPerModule {
	case 1:
		return colorPalette1Bit
	case 2:
		return colorPalette2Bit
	case 4:
		return colorPalette4Bit
	default:
		return ColorPalette
	}
}

func encodeMetadata(errorLevel, bitsPerModule int, isSegment bool) byte {
	metadata := byte(errorLevel & 0x03)
	if bitsPerModule >= 1 && bitsPerModule <= 3 {
		metadata |= byte(bitsPerModule&0x03) << 2
	} else {
		metadata |= byte(bitsPerModule&0x07) << 4
	}
	if isSegment {
		metadata |= 0x80
	}
	return metadata
}

func decodeMetadata(metadata byte) (errorLevel int, bitsPerModule int, isSegment bool, err error) {
	errorLevel = int(metadata & 0x03)
	if errorLevel < 0 || errorLevel > 3 {
		return 0, 0, false, fmt.Errorf("invalid error level %d", errorLevel)
	}

	legacyBits := int((metadata >> 2) & 0x03)
	if legacyBits != 0 {
		bitsPerModule = legacyBits
	} else {
		extendedBits := int((metadata >> 4) & 0x07)
		if extendedBits != 0 {
			bitsPerModule = extendedBits
		} else {
			bitsPerModule = 3 // Backward compatibility with old payloads that stored zero.
		}
	}

	if bitsPerModule < 1 || bitsPerModule > 4 {
		return 0, 0, false, fmt.Errorf("invalid bit density %d in metadata", bitsPerModule)
	}

	isSegment = (metadata & 0x80) != 0
	return errorLevel, bitsPerModule, isSegment, nil
}

func extractPackedBits(data []byte, bitPos, width int) int {
	result := 0
	for i := range width {
		result <<= 1
		pos := bitPos + i
		if pos >= len(data)*8 {
			continue
		}
		byteIdx := pos / 8
		bitOffset := pos % 8
		bit := (data[byteIdx] >> (7 - bitOffset)) & 1
		result |= int(bit)
	}
	return result
}

func extractBytesFromMatrix(matrix [][]int, bitsPerModule int) []byte {
	size := len(matrix)
	var bits []byte
	currentByte := byte(0)
	bitCount := 0

	for y := size - 1; y >= 0; y-- {
		for x := size - 1; x >= 0; x-- {
			if isReservedPos(x, y, size) {
				continue
			}

			colorVal := matrix[y][x]
			for i := bitsPerModule - 1; i >= 0; i-- {
				bit := (colorVal >> i) & 1
				currentByte = (currentByte << 1) | byte(bit)
				bitCount++
				if bitCount == 8 {
					bits = append(bits, currentByte)
					currentByte = 0
					bitCount = 0
				}
			}
		}
	}

	return bits
}

func decodeMatrixWithBits(matrix [][]int, opts *Options, bitsPerModule int) ([]byte, error) {
	bits := extractBytesFromMatrix(matrix, bitsPerModule)
	if len(bits) < 3 {
		return nil, fmt.Errorf("insufficient data: got %d bytes, need at least 3", len(bits))
	}

	metadata := bits[0]
	errorLevel, storedBits, _, metadataErr := decodeMetadata(metadata)
	if metadataErr != nil {
		return nil, metadataErr
	}
	if storedBits != bitsPerModule {
		return nil, fmt.Errorf("bit density mismatch: metadata=%d extracted=%d", storedBits, bitsPerModule)
	}

	dataLen := int(bits[1])<<8 | int(bits[2])
	redundancy := []int{10, 20, 30, 40}[errorLevel]
	totalNeeded := 3 + dataLen
	encodedLen := min((totalNeeded*(100+redundancy)+99)/100, len(bits))
	bits = bits[:encodedLen]

	originalLen := max((len(bits)*100)/(100+redundancy), totalNeeded)
	if originalLen > len(bits) {
		originalLen = len(bits)
	}
	bits = bits[:originalLen]

	if len(bits) < totalNeeded {
		return nil, fmt.Errorf("insufficient data: need %d bytes, have %d", totalNeeded, len(bits))
	}

	bits = bits[3:totalNeeded]
	if len(bits) < 8 {
		return nil, fmt.Errorf("insufficient data for checksum: got %d bytes", len(bits))
	}

	checksumStart := len(bits) - 8
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
	for i := range 8 {
		if checksum[i] != hash[i] {
			return nil, fmt.Errorf("checksum mismatch at byte %d: expected %02x, got %02x", i, hash[i], checksum[i])
		}
	}

	return processed, nil
}
