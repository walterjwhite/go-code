package pubsub

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	publishFunc func(topic string, message []byte) error
	initFunc    func(pctx context.Context) error
	cancelFunc  func()
}

func (m *mockPublisher) Publish(topic string, message []byte) error {
	if m.publishFunc != nil {
		return m.publishFunc(topic, message)
	}
	return nil
}

func (m *mockPublisher) Init(pctx context.Context) error {
	if m.initFunc != nil {
		return m.initFunc(pctx)
	}
	return nil
}

func (m *mockPublisher) Cancel() {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}
}

func TestGetPublisherAndTopic(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	w.setPublisher(mock)

	pub, topic := w.getPublisherAndTopic()
	assert.Equal(t, mock, pub)
	assert.Equal(t, "test-topic", topic)
}

func TestSetPublisher(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{}

	w.setPublisher(mock)

	assert.Equal(t, mock, w.Publisher)
}

func TestGetPublisher(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{}
	w.setPublisher(mock)

	pub := w.getPublisher()
	assert.Equal(t, mock, pub)
}

func TestPubSubWriter_Write(t *testing.T) {
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			assert.Equal(t, "test-topic", topic)
			assert.Equal(t, "hello", string(message))
			return nil
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	n, err := w.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestPubSubWriter_Write_Error(t *testing.T) {
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			return errors.New("publish error")
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	_, err = w.Write([]byte("hello"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message to topic")
}

func TestPubSubWriter_Close(t *testing.T) {
	cancelled := false
	mock := &mockPublisher{
		cancelFunc: func() {
			cancelled = true
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	w.Close()
	assert.True(t, cancelled)
}

func TestPubSubWriter_Init(t *testing.T) {
	initialized := false
	mock := &mockPublisher{
		initFunc: func(pctx context.Context) error {
			initialized = true
			return nil
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)
	assert.True(t, initialized)
}

func TestPubSubWriter_Write_NilPublisher(t *testing.T) {
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), nil)
	assert.NoError(t, err)
	n, err := w.Write([]byte("hello"))
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "publisher not initialized")
}

func TestPubSubWriter_Close_NilPublisher(t *testing.T) {
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), nil)
	assert.NoError(t, err)
	w.Close()
}

func TestPubSubWriter_Init_EmptyTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "",
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "topic name cannot be empty")
}

func TestPubSubWriter_Init_InvalidTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "invalid topic!",
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid topic name")
}

func TestPubSubWriter_Init_TooLongTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "a" + string(make([]byte, 256)),
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid topic name")
}

func TestPubSubWriter_Init_ValidTopicNames(t *testing.T) {
	validNames := []string{
		"simple-topic",
		"topic_with_underscore",
		"Topic123",
		"a",
		"topic-with-multiple_separators-123",
	}

	mock := &mockPublisher{}
	for _, name := range validNames {
		w := &PubSubWriter{
			TopicName: name,
		}
		err := w.Init(context.Background(), mock)
		assert.NoError(t, err, "topic name %q should be valid", name)
	}
}

func TestPubSubWriter_Write_EmptyMessage(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	n, err := w.Write([]byte{})
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "message cannot be empty")
}

func TestPubSubWriter_Write_NilMessage(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	n, err := w.Write(nil)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "message cannot be empty")
}

func TestPubSubWriter_Write_MessageTooLarge(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	largeMessage := make([]byte, maxMessageSize+1)
	n, err := w.Write(largeMessage)
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "exceeds maximum allowed size")
}

func TestPubSubWriter_Write_ValidMessageSizes(t *testing.T) {
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			return nil
		},
	}
	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	n, err := w.Write([]byte("small"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	maxSizeMsg := make([]byte, maxMessageSize)
	n, err = w.Write(maxSizeMsg)
	assert.NoError(t, err)
	assert.Equal(t, maxMessageSize, n)
}

func TestValidateMessage(t *testing.T) {
	err := validateMessage([]byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message cannot be empty")

	err = validateMessage(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message cannot be empty")

	largeMsg := make([]byte, maxMessageSize+1)
	err = validateMessage(largeMsg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed size")

	err = validateMessage([]byte("valid"))
	assert.NoError(t, err)
}

func TestPubSubWriter_Close_MultipleCalls(t *testing.T) {
	cancelCount := 0
	mock := &mockPublisher{
		cancelFunc: func() {
			cancelCount++
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	w.Close()
	w.Close()
	w.Close()

	assert.Equal(t, 3, cancelCount)
}

func TestPubSubWriter_ConcurrentWrite(t *testing.T) {
	var publishCount int
	var mu sync.Mutex
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			mu.Lock()
			publishCount++
			mu.Unlock()
			return nil
		},
	}

	w := &PubSubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	done := make(chan bool)
	for range 10 {
		go func() {
			_, err := w.Write([]byte("concurrent message"))
			assert.NoError(t, err)
			done <- true
		}()
	}

	for range 10 {
		<-done
	}

	mu.Lock()
	count := publishCount
	mu.Unlock()
	assert.Equal(t, 10, count)
}
