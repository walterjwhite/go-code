package zstd

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

type errReader struct{}

func (r *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

type errWriter struct {
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

type errCloser struct {
	*bytes.Buffer
}

func (w *errCloser) Close() error {
	return errors.New("close error")
}

type readErrCloser struct {
	*bytes.Reader
}

func (r *readErrCloser) Close() error {
	return errors.New("close error")
}

type limitedReader struct {
	*bytes.Reader
	n int
}

func (r *limitedReader) Read(p []byte) (n int, err error) {
	if r.n <= 0 {
		return 0, errors.New("limited reader error")
	}
	n, err = r.Reader.Read(p)
	r.n -= n
	return
}

type failingReader struct{}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("failing reader error")
}

func TestStream(t *testing.T) {
	input := []byte("hello, world")
	buf := new(bytes.Buffer)

	err := CompressStream(bytes.NewReader(input), buf)
	if err != nil {
		t.Errorf("CompressStream() error = %v", err)
	}

	compressed := buf.Bytes()

	outputBuf := new(bytes.Buffer)
	err = DecompressStream(bytes.NewReader(compressed), outputBuf)
	if err != nil {
		t.Errorf("DecompressStream() error = %v", err)
	}

	if !bytes.Equal(input, outputBuf.Bytes()) {
		t.Errorf("Stream() = %v, want %v", outputBuf.Bytes(), input)
	}
}

func TestBuffer(t *testing.T) {
	input := []byte("hello, world")
	compressed, err := CompressBuffer(input)
	if err != nil {
		t.Fatalf("CompressBuffer() error = %v", err)
	}

	output, err := DecompressBuffer(compressed)
	if err != nil {
		t.Errorf("DecompressBuffer() error = %v", err)
	}

	if !bytes.Equal(input, output) {
		t.Errorf("Buffer() = %v, want %v", output, input)
	}
}

func TestCompressStreamError(t *testing.T) {
	in := &errReader{}
	var out bytes.Buffer
	err := CompressStream(in, &out)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestCompressStreamErrorOnClose(t *testing.T) {
	in := bytes.NewReader([]byte("hello"))
	out := &errCloser{new(bytes.Buffer)}
	err := CompressStream(in, out)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestCompressStreamErrorOnCopyAndClose(t *testing.T) {
	in := &limitedReader{bytes.NewReader(make([]byte, 1024)), 512}
	out := &errCloser{new(bytes.Buffer)}
	err := CompressStream(in, out)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestCompressStreamCopyAndCloseError(t *testing.T) {
	in := &failingReader{} // Using failingReader to guarantee io.Copy error
	out := &errCloser{new(bytes.Buffer)}

	err := CompressStream(in, out)
	if err == nil {
		t.Fatal("expected an error, but got none")
	}
}

func TestDecompressStreamError(t *testing.T) {
	in := &errReader{}
	var out bytes.Buffer
	err := DecompressStream(in, &out)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestDecompressStreamInvalidData(t *testing.T) {
	in := bytes.NewReader([]byte("invalid data"))
	var out bytes.Buffer
	err := DecompressStream(in, &out)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestDecompressStreamWriteError(t *testing.T) {
	input := []byte("hello, world")
	buf := new(bytes.Buffer)

	err := CompressStream(bytes.NewReader(input), buf)
	if err != nil {
		t.Errorf("CompressStream() error = %v", err)
	}

	compressed := buf.Bytes()

	err = DecompressStream(bytes.NewReader(compressed), &errWriter{})
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestDecompressStreamReadCloser(t *testing.T) {
	input := []byte("hello, world")
	buf := new(bytes.Buffer)

	err := CompressStream(bytes.NewReader(input), buf)
	if err != nil {
		t.Errorf("CompressStream() error = %v", err)
	}

	compressed := buf.Bytes()

	outputBuf := new(bytes.Buffer)
	err = DecompressStream(&readErrCloser{bytes.NewReader(compressed)}, outputBuf)
	if err == nil {
		t.Error("expected an error, but got none")
	}
}

func TestCompressStreamWithContext(t *testing.T) {
	input := []byte("hello, world")
	buf := new(bytes.Buffer)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := CompressStreamWithContext(ctx, bytes.NewReader(input), buf)
	if err != nil {
		t.Errorf("CompressStreamWithContext() error = %v", err)
	}

	compressed := buf.Bytes()

	outputBuf := new(bytes.Buffer)
	err = DecompressStreamWithContext(ctx, bytes.NewReader(compressed), outputBuf)
	if err != nil {
		t.Errorf("DecompressStreamWithContext() error = %v", err)
	}

	if !bytes.Equal(input, outputBuf.Bytes()) {
		t.Errorf("StreamWithContext() = %v, want %v", outputBuf.Bytes(), input)
	}
}

func TestCompressStreamWithContextCancellation(t *testing.T) {
	input := make([]byte, 1024*1024) // 1MB of data
	buf := new(bytes.Buffer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := CompressStreamWithContext(ctx, bytes.NewReader(input), buf)
	if err == nil {
		t.Error("expected context cancellation error, but got none")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestDecompressStreamWithContextCancellation(t *testing.T) {
	input := []byte("hello, world")
	buf := new(bytes.Buffer)

	_ = CompressStream(bytes.NewReader(input), buf)
	compressed := buf.Bytes()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	outputBuf := new(bytes.Buffer)
	err := DecompressStreamWithContext(ctx, bytes.NewReader(compressed), outputBuf)
	if err == nil {
		t.Error("expected context cancellation error, but got none")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestDecompressBufferInvalidData(t *testing.T) {
	invalid := []byte("this is not compressed data")
	_, err := DecompressBuffer(invalid)
	if err == nil {
		t.Error("expected error for invalid compressed data, got nil")
	}
}

func TestCompressBufferNilInput(t *testing.T) {
	_, err := CompressBuffer(nil)
	if err != nil {
		t.Errorf("CompressBuffer(nil) should not error, got %v", err)
	}
}
