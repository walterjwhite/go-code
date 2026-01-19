package zstd

import (
	"bytes"
	"errors"
	"strings" // Added for error string checks
	"testing"
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
	compressed := CompressBuffer(input)

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
	in := &failingReader{}               // Using failingReader to guarantee io.Copy error
	out := &errCloser{new(bytes.Buffer)} // enc.Close will error

	err := CompressStream(in, out)
	if err == nil {
		t.Fatal("expected an error, but got none")
	}

	expectedErrPart := "Second error" // from fmt.Errorf("%w; Second error", err2)
	if !strings.Contains(err.Error(), expectedErrPart) {
		t.Errorf("expected error message to contain '%s', got '%v'", expectedErrPart, err)
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
