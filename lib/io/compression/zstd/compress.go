package zstd

import (
	"context"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

const DefaultMaxEncodedSize = 1024 * 1024 * 1024 // 1 GB

var encoderOptions = []zstd.EOption{
	zstd.WithEncoderConcurrency(0), // Use GOMAXPROCS
	zstd.WithEncoderCRC(true),      // Enable CRC for integrity checking
}

func getEncoder(out io.Writer) (*zstd.Encoder, error) {
	return zstd.NewWriter(out, encoderOptions...)
}

func CompressStream(in io.Reader, out io.Writer) (err error) {
	return CompressStreamWithContext(context.Background(), in, out)
}

func CompressStreamWithContext(ctx context.Context, in io.Reader, out io.Writer) (err error) {
	if closer, ok := out.(io.Closer); ok {
		defer func() {
			if cerr := closer.Close(); err == nil {
				err = cerr
			}
		}()
	}

	enc, err := getEncoder(out)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(enc, in)
		done <- copyErr
	}()

	select {
	case <-ctx.Done():
		if cerr := enc.Close(); cerr != nil && err == nil {
			err = cerr
		}
		return ctx.Err()
	case copyErr := <-done:
		if cerr := enc.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if copyErr != nil {
			return copyErr
		}
	}

	return nil
}

func CompressBuffer(src []byte) ([]byte, error) {
	enc, err := zstd.NewWriter(nil, encoderOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}
	defer func() { _ = enc.Close() }()

	result := enc.EncodeAll(src, make([]byte, 0, len(src)))
	return result, nil
}
