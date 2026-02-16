package zstd

import (
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io"
)

func CompressStream(in io.Reader, out io.Writer) (err error) {
	if closer, ok := out.(io.Closer); ok {
		defer func() {
			if cerr := closer.Close(); err == nil {
				err = cerr
			}
		}()
	}
	enc, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	_, err = io.Copy(enc, in)
	if err != nil {
		err2 := enc.Close()
		if err2 != nil {
			err = fmt.Errorf("%w; Second error", err2)
		}

		return err
	}
	return enc.Close()
}

var encoder, _ = zstd.NewWriter(nil)

func CompressBuffer(src []byte) []byte {
	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}
