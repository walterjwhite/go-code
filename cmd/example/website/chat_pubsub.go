package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
	googlepubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type pubsubChatService struct {
	conf                 *googlepubsub.Conf
	publishTopic         string
	responseTopic        string
	responseSubscription string
	limiter              *chatlib.IPRateLimiter
	publishFn            func(topic string, payload []byte) error
	responses            chan chatlib.ServerMessage
	errors               chan error
}

type clientMessageEnvelope struct {
	IP            string            `json:"ip,omitempty"`
	IPAddress     string            `json:"ip_address,omitempty"`
	Message       string            `json:"message"`
	Channel       string            `json:"channel,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SentAt        time.Time         `json:"sent_at"`
}

type serverMessageEnvelope struct {
	IP            string            `json:"ip,omitempty"`
	IPAddress     string            `json:"ip_address,omitempty"`
	Message       string            `json:"message"`
	Channel       string            `json:"channel,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SentAt        time.Time         `json:"sent_at"`
}

var websiteChatRateLimit = chatlib.RateLimitConfig{
	Messages: 20,
	Window:   time.Minute,
}

func newPubsubChatService(
	conf *googlepubsub.Conf,
	publishTopic, responseTopic, responseSubscription string,
) *pubsubChatService {
	limiter := chatlib.NewIPRateLimiter(websiteChatRateLimit)
	limiter.Start(context.Background(), websiteChatRateLimit.Window)

	return &pubsubChatService{
		conf:                 conf,
		publishTopic:         publishTopic,
		responseTopic:        responseTopic,
		responseSubscription: responseSubscription,
		limiter:              limiter,
		publishFn:            conf.Publish,
		responses:            make(chan chatlib.ServerMessage, chatlib.DefaultResponseBuffer),
		errors:               make(chan error, 8),
	}
}

func (s *pubsubChatService) Send(ctx context.Context, msg chatlib.ClientMessage) error {
	if ctx == nil {
		return chatlib.ErrNilContext
	}
	if msg.IP == "" {
		return chatlib.ErrEmptyIP
	}
	if msg.Message == "" {
		return chatlib.ErrEmptyMessage
	}
	if !s.limiter.Allow(msg.IP) {
		return chatlib.ErrRateLimited
	}

	payload, err := marshalClientMessage(msg)
	if err != nil {
		return fmt.Errorf("serialize client message: %w", err)
	}

	if err := s.publishFn(s.publishTopic, payload); err != nil {
		return fmt.Errorf("publish client message: %w", err)
	}

	log.Debug().
		Str("ip", msg.IP).
		Str("correlation_id", msg.CorrelationID).
		Str("topic", s.publishTopic).
		Msg("pubsub: published client message")

	return nil
}

func (s *pubsubChatService) Start(ctx context.Context) {
	sub := &websiteSubscriber{
		responses: s.responses,
		errors:    s.errors,
	}
	go func() {
		log.Info().
			Str("topic", s.responseTopic).
			Str("subscription", s.responseSubscription).
			Msg("pubsub: starting response subscription")

		s.conf.Subscribe(s.responseSubscription, sub)

		log.Warn().Msg("pubsub: response subscription ended")
	}()
}

func (s *pubsubChatService) Responses() <-chan chatlib.ServerMessage { return s.responses }
func (s *pubsubChatService) Errors() <-chan error                    { return s.errors }


type websiteSubscriber struct {
	responses chan chatlib.ServerMessage
	errors    chan error
}

func (ws *websiteSubscriber) MessageDeserialized(data []byte) {
	msg, err := unmarshalServerMessage(data)
	if err != nil {
		ws.pushError(fmt.Errorf("deserialize server message: %w", err))
		return
	}

	log.Debug().
		Str("correlation_id", msg.CorrelationID).
		Str("ip", msg.Metadata["ip"]).
		Str("message", msg.Message).
		Msg("pubsub: received server message")

	select {
	case ws.responses <- msg:
	default:
		log.Warn().
			Str("correlation_id", msg.CorrelationID).
			Msg("pubsub: responses channel full — dropping message")
	}
}

func marshalClientMessage(msg chatlib.ClientMessage) ([]byte, error) {
	if msg.SentAt.IsZero() {
		msg.SentAt = time.Now().UTC()
	}

	if msg.Metadata == nil {
		msg.Metadata = make(map[string]string)
	}
	msg.Metadata["ip"] = msg.IP

	return json.Marshal(clientMessageEnvelope{
		IP:            msg.IP,
		IPAddress:     msg.IP,
		Message:       msg.Message,
		Channel:       msg.Channel,
		CorrelationID: msg.CorrelationID,
		Metadata:      msg.Metadata,
		SentAt:        msg.SentAt,
	})
}

func unmarshalServerMessage(data []byte) (chatlib.ServerMessage, error) {
	var envelope serverMessageEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return chatlib.ServerMessage{}, err
	}

	msg := chatlib.ServerMessage{
		Channel:       envelope.Channel,
		CorrelationID: envelope.CorrelationID,
		Message:       envelope.Message,
		Metadata:      envelope.Metadata,
		SentAt:        envelope.SentAt,
	}
	if msg.Metadata == nil {
		msg.Metadata = make(map[string]string)
	}
	if msg.Metadata["ip"] == "" {
		switch {
		case envelope.IP != "":
			msg.Metadata["ip"] = envelope.IP
		case envelope.IPAddress != "":
			msg.Metadata["ip"] = envelope.IPAddress
		}
	}

	return msg, nil
}

func (ws *websiteSubscriber) MessageParseError(err error) {
	ws.pushError(fmt.Errorf("pubsub message parse error: %w", err))
}

func (ws *websiteSubscriber) pushError(err error) {
	select {
	case ws.errors <- err:
	default:
		log.Error().Err(err).Msg("pubsub: errors channel full")
	}
}
