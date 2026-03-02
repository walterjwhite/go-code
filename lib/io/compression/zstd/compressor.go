package zstd

import (
	"github.com/walterjwhite/go-code/lib/io/compression"
)

type ZstdCompressor struct{}

func NewCompressor() compression.Compressor {
	return &ZstdCompressor{}
}

func (c *ZstdCompressor) Compress(data []byte) ([]byte, error) {
	return CompressBuffer(data), nil
}

func (c *ZstdCompressor) Decompress(data []byte) ([]byte, error) {
	return DecompressBuffer(data)
}
