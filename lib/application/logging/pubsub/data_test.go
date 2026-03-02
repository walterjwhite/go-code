package pubsub

import (
	"context"
	"errors"
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

func TestPubsubWriter_Write(t *testing.T) {
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			assert.Equal(t, "test-topic", topic)
			assert.Equal(t, "hello", string(message))
			return nil
		},
	}

	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	n, err := w.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestPubsubWriter_Write_Error(t *testing.T) {
	mock := &mockPublisher{
		publishFunc: func(topic string, message []byte) error {
			return errors.New("publish error")
		},
	}

	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	_, err = w.Write([]byte("hello"))
	assert.Error(t, err)
}

func TestPubsubWriter_Close(t *testing.T) {
	cancelled := false
	mock := &mockPublisher{
		cancelFunc: func() {
			cancelled = true
		},
	}

	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)

	w.Close()
	assert.True(t, cancelled)
}

func TestPubsubWriter_Init(t *testing.T) {
	initialized := false
	mock := &mockPublisher{
		initFunc: func(pctx context.Context) error {
			initialized = true
			return nil
		},
	}

	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), mock)
	assert.NoError(t, err)
	assert.True(t, initialized)
}

func TestPubsubWriter_Write_NilPublisher(t *testing.T) {
	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), nil)
	assert.NoError(t, err)
	n, err := w.Write([]byte("hello"))
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "publisher not initialized")
}

func TestPubsubWriter_Close_NilPublisher(t *testing.T) {
	w := &PubsubWriter{
		TopicName: "test-topic",
	}
	err := w.Init(context.Background(), nil)
	assert.NoError(t, err)
	w.Close()
}

func TestPubsubWriter_Init_EmptyTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubsubWriter{
		TopicName: "",
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "topic name cannot be empty")
}

func TestPubsubWriter_Init_InvalidTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubsubWriter{
		TopicName: "invalid topic!",
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid topic name")
}

func TestPubsubWriter_Init_TooLongTopicName(t *testing.T) {
	mock := &mockPublisher{}
	w := &PubsubWriter{
		TopicName: "a" + string(make([]byte, 256)),
	}
	err := w.Init(context.Background(), mock)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid topic name")
}

func TestPubsubWriter_Init_ValidTopicNames(t *testing.T) {
	validNames := []string{
		"simple-topic",
		"topic_with_underscore",
		"Topic123",
		"a",
		"topic-with-multiple_separators-123",
	}

	mock := &mockPublisher{}
	for _, name := range validNames {
		w := &PubsubWriter{
			TopicName: name,
		}
		err := w.Init(context.Background(), mock)
		assert.NoError(t, err, "topic name %q should be valid", name)
	}
}
