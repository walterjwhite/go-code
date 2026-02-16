package pubsub

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	publishFunc func(topic string, message []byte) error
	initFunc    func(pctx context.Context)
	cancelFunc  func()
}

func (m *mockPublisher) Publish(topic string, message []byte) error {
	if m.publishFunc != nil {
		return m.publishFunc(topic, message)
	}
	return nil
}

func (m *mockPublisher) Init(pctx context.Context) {
	if m.initFunc != nil {
		m.initFunc(pctx)
	}
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
	w.Init(context.Background(), mock)

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
	w.Init(context.Background(), mock)

	_, err := w.Write([]byte("hello"))
	assert.Error(t, err)
}

func TestPubsubWriter_Close(t *testing.T) {
	cancelled := false
	mock := &mockPublisher{
		cancelFunc: func() {
			cancelled = true
		},
	}

	w := &PubsubWriter{}
	w.Init(context.Background(), mock)

	w.Close()
	assert.True(t, cancelled)
}

func TestPubsubWriter_Init(t *testing.T) {
	initialized := false
	mock := &mockPublisher{
		initFunc: func(pctx context.Context) {
			initialized = true
		},
	}

	w := &PubsubWriter{}
	w.Init(context.Background(), mock)
	assert.True(t, initialized)
}

func TestPubsubWriter_Write_NilPublisher(t *testing.T) {
	w := &PubsubWriter{}
	w.Init(context.Background(), nil)
	n, err := w.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestPubsubWriter_Close_NilPublisher(t *testing.T) {
	w := &PubsubWriter{}
	w.Init(context.Background(), nil)
	w.Close()
}
