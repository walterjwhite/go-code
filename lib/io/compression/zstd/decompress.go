package zstd

import (
	"github.com/klauspost/compress/zstd"
	"io"
)

func DecompressStream(in io.Reader, out io.Writer) error {
	d, err := zstd.NewReader(in)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(out, d)
	return err
}

var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))

func DecompressBuffer(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}
