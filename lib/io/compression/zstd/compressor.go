package zstd

import (
	"context"

	"github.com/walterjwhite/go-code/lib/io/compression"
)

type ZstdCompressor struct{}

func NewCompressor() compression.Compressor {
	return &ZstdCompressor{}
}

func (c *ZstdCompressor) Compress(data []byte) ([]byte, error) {
	return CompressBuffer(data)
}

func (c *ZstdCompressor) Decompress(data []byte) ([]byte, error) {
	return DecompressBuffer(data)
}

type ContextCompressor interface {
	compression.Compressor
	CompressWithContext(ctx context.Context, data []byte) ([]byte, error)
	DecompressWithContext(ctx context.Context, data []byte) ([]byte, error)
}

type ZstdContextCompressor struct{}

func NewContextCompressor() ContextCompressor {
	return &ZstdContextCompressor{}
}

func (c *ZstdContextCompressor) Compress(data []byte) ([]byte, error) {
	return c.CompressWithContext(context.Background(), data)
}

func (c *ZstdContextCompressor) Decompress(data []byte) ([]byte, error) {
	return c.DecompressWithContext(context.Background(), data)
}

func (c *ZstdContextCompressor) CompressWithContext(ctx context.Context, data []byte) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return CompressBuffer(data)
}

func (c *ZstdContextCompressor) DecompressWithContext(ctx context.Context, data []byte) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return DecompressBuffer(data)
}
