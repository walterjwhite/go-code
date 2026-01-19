package pipe

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFlusher struct {
	mock.Mock
	mu       sync.Mutex
	flushedC chan []byte
}

func (m *MockFlusher) Flush(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(data)

	if m.flushedC != nil {
		m.flushedC <- data
	}
	return args.Error(0)
}

func TestReader_Flush(t *testing.T) {
	t.Run("empty buffer does nothing", func(t *testing.T) {
		r := &Reader{}
		err := r.Flush()
		assert.NoError(t, err)
	})

	t.Run("nil flusher does not error", func(t *testing.T) {
		r := &Reader{}
		r.buf.Write([]byte("some data"))
		err := r.Flush()
		assert.NoError(t, err)
		assert.Equal(t, 0, r.buf.Len(), "buffer should be empty after flush")
	})

	t.Run("flushes buffer contents", func(t *testing.T) {
		mockFlusher := new(MockFlusher)
		r := &Reader{
			Flusher: mockFlusher,
		}
		testData := []byte("hello")
		r.buf.Write(testData)

		mockFlusher.On("Flush", testData).Return(nil).Once()

		err := r.Flush()
		assert.NoError(t, err)

		mockFlusher.AssertExpectations(t)
		assert.Equal(t, 0, r.buf.Len(), "buffer should be empty after flush")
	})
}

func TestReader_Start(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pipe-test")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	pipePath := filepath.Join(tempDir, "test.pipe")
	err = syscall.Mkfifo(pipePath, 0666)
	assert.NoError(t, err)

	mockFlusher := &MockFlusher{
		flushedC: make(chan []byte, 1),
	}
	mockFlusher.On("Flush", mock.Anything).Return(nil)

	r := &Reader{
		PipePath:  pipePath,
		Threshold: 10,
		Flusher:   mockFlusher,
	}
	r.Start()

	time.Sleep(100 * time.Millisecond)

	f, err := os.OpenFile(pipePath, os.O_WRONLY, 0)
	assert.NoError(t, err)

	_, err = f.WriteString("hello")
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // allow reader to process
	mockFlusher.AssertNotCalled(t, "Flush", mock.Anything)
	assert.Equal(t, 5, r.buf.Len())

	_, err = f.WriteString(" world")
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "hello world", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for flush")
	}

	assert.Equal(t, 0, r.buf.Len())

	_, err = f.WriteString("manual")
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond) // allow reader to process
	assert.Equal(t, 6, r.buf.Len())

	err = r.Flush()
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "manual", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for manual flush")
	}
	assert.Equal(t, 0, r.buf.Len())

	err = f.Close()
	assert.NoError(t, err)
	time.Sleep(600 * time.Millisecond) // allow reader to detect close and reopen

	f, err = os.OpenFile(pipePath, os.O_WRONLY, 0)
	assert.NoError(t, err)
	defer func() {
		err := f.Close()
		assert.NoError(t, err)
	}()

	_, err = f.WriteString("after reopen and more")
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "after reopen and more", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for flush after reopen")
	}

	mockFlusher.AssertNumberOfCalls(t, "Flush", 3)
}

func TestReader_Start_FlusherError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pipe-test-flusher-error")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	pipePath := filepath.Join(tempDir, "test-flusher-error.pipe")
	err = syscall.Mkfifo(pipePath, 0666)
	assert.NoError(t, err)

	mockFlusher := &MockFlusher{
		flushedC: make(chan []byte, 1),
	}
	flushErr := assert.AnError
	mockFlusher.On("Flush", mock.Anything).Return(flushErr).Once() // Only expect one error flush

	r := &Reader{
		PipePath:  pipePath,
		Threshold: 5, // Small threshold to trigger flush quickly
		Flusher:   mockFlusher,
	}
	r.Start()

	time.Sleep(100 * time.Millisecond)

	f, err := os.OpenFile(pipePath, os.O_WRONLY, 0)
	assert.NoError(t, err)
	defer func() {
		err := f.Close()
		assert.NoError(t, err)
	}()

	_, err = f.WriteString("trigger flush error") // This write should exceed threshold
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)

	select {
	case <-mockFlusher.flushedC:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for error flush data to be sent to channel")
	}

	mockFlusher.AssertCalled(t, "Flush", []byte("trigger flush error"))
	mockFlusher.AssertExpectations(t)
	assert.Equal(t, 0, r.buf.Len(), "buffer should be cleared even if flush fails")

	mockFlusher.On("Flush", mock.Anything).Return(nil).Once()
	_, err = f.WriteString("recovery")
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "recovery", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for recovery flush")
	}
}

func TestReader_Start_OpenFileError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pipe-test-open-error")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	pipePath := filepath.Join(tempDir, "non-existent.pipe") // Initially, this pipe does not exist

	mockFlusher := &MockFlusher{
		flushedC: make(chan []byte, 1),
	}
	mockFlusher.On("Flush", mock.Anything).Return(nil)

	r := &Reader{
		PipePath:  pipePath,
		Threshold: 10,
		Flusher:   mockFlusher,
	}
	r.Start()

	time.Sleep(1 * time.Second)

	err = syscall.Mkfifo(pipePath, 0666)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	f, err := os.OpenFile(pipePath, os.O_WRONLY, 0)
	assert.NoError(t, err)
	defer func() {
		err := f.Close()
		assert.NoError(t, err)
	}()

	_, err = f.WriteString("recovery data")
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "recovery data", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for flush after recovery")
	}

	mockFlusher.AssertCalled(t, "Flush", []byte("recovery data"))
}

func TestReader_Start_ThresholdDisabled(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pipe-test-threshold-disabled")
	assert.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDir)
		assert.NoError(t, err)
	}()

	pipePath := filepath.Join(tempDir, "test-threshold-disabled.pipe")
	err = syscall.Mkfifo(pipePath, 0666)
	assert.NoError(t, err)

	mockFlusher := &MockFlusher{
		flushedC: make(chan []byte, 1),
	}
	mockFlusher.On("Flush", mock.Anything).Return(nil)

	r := &Reader{
		PipePath:  pipePath,
		Threshold: 0, // Threshold is 0, so automatic flush should be disabled
		Flusher:   mockFlusher,
	}
	r.Start()

	time.Sleep(100 * time.Millisecond)

	f, err := os.OpenFile(pipePath, os.O_WRONLY, 0)
	assert.NoError(t, err)
	defer func() {
		err := f.Close()
		assert.NoError(t, err)
	}()

	_, err = f.WriteString("this is more than enough data to trigger a flush normally")
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond) // Give reader time to process

	mockFlusher.AssertNotCalled(t, "Flush", mock.Anything)
	assert.Greater(t, r.buf.Len(), 0, "buffer should contain data as no auto-flush occurred")

	err = r.Flush()
	assert.NoError(t, err)

	select {
	case flushed := <-mockFlusher.flushedC:
		assert.Equal(t, "this is more than enough data to trigger a flush normally", string(flushed))
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for manual flush")
	}

	mockFlusher.AssertCalled(t, "Flush", []byte("this is more than enough data to trigger a flush normally"))
	assert.Equal(t, 0, r.buf.Len(), "buffer should be empty after manual flush")
}

func TestReader_ConcurrentFlush(t *testing.T) {
	mockFlusher := &MockFlusher{}
	mockFlusher.On("Flush", mock.Anything).Return(nil)

	r := &Reader{
		Flusher: mockFlusher,
	}

	data := []byte("some concurrent data")
	r.mu.Lock()
	r.buf.Write(data)
	r.mu.Unlock()

	var wg sync.WaitGroup
	numConcurrentCalls := 10

	for i := 0; i < numConcurrentCalls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.Flush()
		}()
	}

	wg.Wait()

	mockFlusher.AssertCalled(t, "Flush", data)
	mockFlusher.AssertNumberOfCalls(t, "Flush", 1)
	assert.Equal(t, 0, r.buf.Len(), "buffer should be empty after concurrent flushes")
}
