package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	googlepubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type mockPubSub struct {
	publishedTopic   string
	publishedPayload []byte
}

func (m *mockPubSub) Publish(topicName string, message []byte) error {
	m.publishedTopic = topicName
	m.publishedPayload = append([]byte(nil), message...)
	return nil
}

func (m *mockPubSub) Subscribe(subscriptionName string, ms googlepubsub.MessageSubscriber) {
}

func TestNewWithQuery_Validation(t *testing.T) {
	_, err := NewWithQuery(Config{}, nil)
	assert.Error(t, err)

	_, err = NewWithQuery(Config{
		PubSub:              &mockPubSub{},
		RequestSubscription: "sub",
		ResponseTopic:       "resp",
	}, nil)
	assert.Error(t, err)
}

func TestService_MessageDeserialized(t *testing.T) {
	pubsub := &mockPubSub{}
	svc, err := NewWithQuery(Config{
		PubSub:              pubsub,
		RequestSubscription: "requests-sub",
		ResponseTopic:       "responses",
		QueryTimeout:        time.Second,
	}, func(ctx context.Context, question string) (string, error) {
		assert.Equal(t, "what is in the docs?", question)
		return "rag answer", nil
	})
	assert.NoError(t, err)

	payload, err := json.Marshal(Message{
		IPAddress: "192.0.2.1",
		Message:   "what is in the docs?",
	})
	assert.NoError(t, err)

	svc.MessageDeserialized(payload)

	assert.Equal(t, "responses", pubsub.publishedTopic)

	var response Message
	err = json.Unmarshal(pubsub.publishedPayload, &response)
	assert.NoError(t, err)
	assert.Equal(t, "192.0.2.1", response.IPAddress)
	assert.Equal(t, "rag answer", response.Message)
}

func TestService_MessageDeserialized_QueryErrorPublishesErrorText(t *testing.T) {
	pubsub := &mockPubSub{}
	svc, err := NewWithQuery(Config{
		PubSub:              pubsub,
		RequestSubscription: "requests-sub",
		ResponseTopic:       "responses",
		QueryTimeout:        time.Second,
	}, func(ctx context.Context, question string) (string, error) {
		return "", errors.New("backend unavailable")
	})
	assert.NoError(t, err)

	payload, err := json.Marshal(Message{
		IPAddress: "198.51.100.2",
		Message:   "question",
	})
	assert.NoError(t, err)

	svc.MessageDeserialized(payload)

	var response Message
	err = json.Unmarshal(pubsub.publishedPayload, &response)
	assert.NoError(t, err)
	assert.Equal(t, "198.51.100.2", response.IPAddress)
	assert.Equal(t, "query error: backend unavailable", response.Message)
}
