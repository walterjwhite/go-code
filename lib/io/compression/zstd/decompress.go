package zstd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

const DefaultMaxDecompressedSize = 1024 * 1024 * 1024 // 1 GB

var ErrDecompressionLimitExceeded = errors.New("decompression limit exceeded: potential decompression bomb")

var decoderOptions = []zstd.DOption{
	zstd.WithDecoderConcurrency(0), // Use GOMAXPROCS
	zstd.WithDecoderMaxMemory(DefaultMaxDecompressedSize),
}

func getDecoder(in io.Reader) (*zstd.Decoder, error) {
	return zstd.NewReader(in, decoderOptions...)
}

func DecompressStream(in io.Reader, out io.Writer) error {
	return DecompressStreamWithContext(context.Background(), in, out)
}

func DecompressStreamWithContext(ctx context.Context, in io.Reader, out io.Writer) (err error) {
	if closer, ok := in.(io.Closer); ok {
		defer func() {
			if cerr := closer.Close(); err == nil {
				err = cerr
			}
		}()
	}

	d, err := getDecoder(in)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}
	defer d.Close()

	limitedOut := &limitingWriter{w: out, limit: DefaultMaxDecompressedSize}

	done := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(limitedOut, d)
		done <- copyErr
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case copyErr := <-done:
		if copyErr != nil {
			if errors.Is(copyErr, errWriteLimitExceeded) {
				return ErrDecompressionLimitExceeded
			}
			return copyErr
		}
	}

	return nil
}

type limitingWriter struct {
	w     io.Writer
	limit int64
	wrote int64
}

var errWriteLimitExceeded = errors.New("write limit exceeded")

func (lw *limitingWriter) Write(p []byte) (int, error) {
	if lw.wrote+int64(len(p)) > lw.limit {
		return 0, errWriteLimitExceeded
	}
	n, err := lw.w.Write(p)
	lw.wrote += int64(n)
	return n, err
}

func DecompressBuffer(src []byte) ([]byte, error) {
	d, err := zstd.NewReader(nil, decoderOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}
	defer d.Close()

	result, err := d.DecodeAll(src, nil)
	if err != nil {
		return nil, err
	}

	if int64(len(result)) > DefaultMaxDecompressedSize {
		return nil, ErrDecompressionLimitExceeded
	}

	return result, nil
}
