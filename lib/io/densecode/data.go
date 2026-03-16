package densecode

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

type Configuration struct {
	data []byte
	size int

	ErrorLevel     int
	ModuleSize     int
	BitsPerModule  int
	Profile        string
	MaxSegmentSize int

	Compressor compression.Compressor
	Encryptor  encryption.Encryptor
}

func WithDefaults() *Configuration {
	return &Configuration{
		ErrorLevel:     1,
		ModuleSize:     10,
		BitsPerModule:  3,
		Profile:        "camera-safe",
		MaxSegmentSize: 32 * 1024,
	}
}

func (c *Configuration) PostLoad(ctx context.Context) []error {
	var validationErrors []error

	if c.ErrorLevel < 0 || c.ErrorLevel > 3 {
		validationErrors = append(validationErrors, errors.New("error level must be in range 0..3"))
	}
	if c.ModuleSize < 1 {
		validationErrors = append(validationErrors, errors.New("module size must be >= 1"))
	}
	if c.MaxSegmentSize < 256 {
		validationErrors = append(validationErrors, errors.New("max segment size must be >= 256"))
	}

	c.Profile = strings.ToLower(strings.TrimSpace(c.Profile))

	switch c.Profile {
	case "camera-safe":
		c.BitsPerModule = 2
	case "balanced":
		c.BitsPerModule = 3
	case "max":
		c.BitsPerModule = 4
	default:
		validationErrors = append(validationErrors, fmt.Errorf("invalid profile %q (allowed: camera-safe|balanced|max)", c.Profile))
	}

	return validationErrors
}

func (c *Configuration) initialize() {
	if c.ErrorLevel == 0 {
		c.ErrorLevel = 1
	}
	if c.ModuleSize == 0 {
		c.ModuleSize = 10
	}
	if c.BitsPerModule == 0 {
		c.BitsPerModule = 3
	}
	if c.MaxSegmentSize == 0 {
		c.MaxSegmentSize = 32 * 1024
	}
	if c.Profile == "" {
		c.Profile = "camera-safe"
	}
}
