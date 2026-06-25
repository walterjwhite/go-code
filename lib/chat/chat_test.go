package chat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/walterjwhite/go-code/lib/io/serialization"
	google_pubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type mockPubSub struct {
	publishedTopic   string
	publishedMessage []byte
	publishErr       error
	subscriber       google_pubsub.MessageSubscriber
}

func (m *mockPubSub) Publish(topic string, message []byte) error {
	m.publishedTopic = topic
	m.publishedMessage = append([]byte(nil), message...)
	return m.publishErr
}

func (m *mockPubSub) Subscribe(_ string, ms google_pubsub.MessageSubscriber) {
	m.subscriber = ms
}

func TestNewClient_ValidateConfig(t *testing.T) {
	serializer := serialization.NewJSONSerializer()
	_, err := NewClient(Config{}, &mockPubSub{}, serializer)
	require.ErrorIs(t, err, ErrPublishTopicRequired)
}

func TestClientSend_PublishesSerializedMessage(t *testing.T) {
	ps := &mockPubSub{}
	serializer := serialization.NewJSONSerializer()

	client, err := NewClient(Config{
		PublishTopic:         "chat-inbound",
		ResponseTopic:        "chat-responses",
		ResponseSubscription: "chat-responses-sub",
		RateLimit: RateLimitConfig{
			Messages: 2,
			Window:   time.Minute,
		},
	}, ps, serializer)
	require.NoError(t, err)

	msg := ClientMessage{
		IP:            "203.0.113.7",
		Message:       "hello",
		Channel:       "room-1",
		CorrelationID: "abc123",
	}

	err = client.Send(context.Background(), msg)
	require.NoError(t, err)
	require.Equal(t, "chat-inbound", ps.publishedTopic)

	var decoded ClientMessage
	require.NoError(t, serializer.Deserialize(ps.publishedMessage, &decoded))
	require.Equal(t, msg.IP, decoded.IP)
	require.Equal(t, msg.Message, decoded.Message)
	require.False(t, decoded.SentAt.IsZero())
}

func TestClientSend_RateLimitedPerIP(t *testing.T) {
	ps := &mockPubSub{}
	serializer := serialization.NewJSONSerializer()

	client, err := NewClient(Config{
		PublishTopic:         "chat-inbound",
		ResponseTopic:        "chat-responses",
		ResponseSubscription: "chat-responses-sub",
		RateLimit: RateLimitConfig{
			Messages: 1,
			Window:   time.Minute,
		},
	}, ps, serializer)
	require.NoError(t, err)

	err = client.Send(context.Background(), ClientMessage{IP: "198.51.100.10", Message: "first"})
	require.NoError(t, err)

	err = client.Send(context.Background(), ClientMessage{IP: "198.51.100.10", Message: "second"})
	require.ErrorIs(t, err, ErrRateLimited)

	err = client.Send(context.Background(), ClientMessage{IP: "198.51.100.11", Message: "other ip"})
	require.NoError(t, err)
}

func TestClientStart_DeliversResponses(t *testing.T) {
	ps := &mockPubSub{}
	serializer := serialization.NewJSONSerializer()

	client, err := NewClient(Config{
		PublishTopic:         "chat-inbound",
		ResponseTopic:        "chat-responses",
		ResponseSubscription: "chat-responses-sub",
		ResponseBuffer:       1,
		RateLimit: RateLimitConfig{
			Messages: 1,
			Window:   time.Minute,
		},
	}, ps, serializer)
	require.NoError(t, err)

	ctx := t.Context()

	client.Start(ctx)
	require.Eventually(t, func() bool {
		return ps.subscriber != nil
	}, time.Second, 10*time.Millisecond)

	payload, err := serializer.Serialize(ServerMessage{
		Channel:       "room-1",
		CorrelationID: "abc123",
		Message:       "reply",
		SentAt:        time.Now().UTC(),
	})
	require.NoError(t, err)

	ps.subscriber.MessageDeserialized(payload)

	select {
	case response := <-client.Responses():
		require.Equal(t, "reply", response.Message)
		require.Equal(t, "abc123", response.CorrelationID)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for response")
	}
}

func TestClientStart_ReportsDeserializeErrors(t *testing.T) {
	ps := &mockPubSub{}
	serializer := serialization.NewJSONSerializer()

	client, err := NewClient(Config{
		PublishTopic:         "chat-inbound",
		ResponseTopic:        "chat-responses",
		ResponseSubscription: "chat-responses-sub",
		ResponseBuffer:       1,
		RateLimit: RateLimitConfig{
			Messages: 1,
			Window:   time.Minute,
		},
	}, ps, serializer)
	require.NoError(t, err)

	ctx := t.Context()

	client.Start(ctx)
	require.Eventually(t, func() bool {
		return ps.subscriber != nil
	}, time.Second, 10*time.Millisecond)

	ps.subscriber.MessageDeserialized([]byte("not-json"))

	select {
	case err := <-client.Errors():
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for error")
	}
}

func TestClientSend_PropagatesPublishError(t *testing.T) {
	ps := &mockPubSub{publishErr: errors.New("boom")}
	serializer := serialization.NewJSONSerializer()

	client, err := NewClient(Config{
		PublishTopic:         "chat-inbound",
		ResponseTopic:        "chat-responses",
		ResponseSubscription: "chat-responses-sub",
		RateLimit: RateLimitConfig{
			Messages: 1,
			Window:   time.Minute,
		},
	}, ps, serializer)
	require.NoError(t, err)

	err = client.Send(context.Background(), ClientMessage{IP: "203.0.113.5", Message: "hello"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "publish chat message")
}
