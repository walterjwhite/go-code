package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/ai/rag"
	"github.com/walterjwhite/go-code/lib/application/logging"
	googlepubsub "github.com/walterjwhite/go-code/lib/net/google"
)

type Message struct {
	IPAddress string `json:"ip_address"`
	Message   string `json:"message"`
}

type PubSub interface {
	Publish(topicName string, message []byte) error
	Subscribe(subscriptionName string, ms googlepubsub.MessageSubscriber)
}

type Config struct {
	Context             context.Context
	PubSub              PubSub
	RequestTopic        string
	RequestSubscription string
	ResponseTopic       string
	QueryTimeout        time.Duration
}

type Service struct {
	cfg   Config
	query func(ctx context.Context, question string) (string, error)
}

func New(cfg Config, engine *rag.Engine) (*Service, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine is required")
	}
	return NewWithQuery(cfg, func(ctx context.Context, question string) (string, error) {
		answer, _, err := engine.QueryText(ctx, question)
		return answer, err
	})
}

func NewWithQuery(cfg Config, query func(ctx context.Context, question string) (string, error)) (*Service, error) {
	if cfg.PubSub == nil {
		return nil, fmt.Errorf("pubsub is required")
	}
	if cfg.RequestSubscription == "" {
		return nil, fmt.Errorf("request subscription is required")
	}
	if cfg.ResponseTopic == "" {
		return nil, fmt.Errorf("response topic is required")
	}
	if query == nil {
		return nil, fmt.Errorf("query function is required")
	}
	if cfg.Context == nil {
		cfg.Context = context.Background()
	}
	if cfg.QueryTimeout <= 0 {
		cfg.QueryTimeout = 1 * time.Minute
	}
	return &Service{cfg: cfg, query: query}, nil
}

func (s *Service) Run() {
	s.cfg.PubSub.Subscribe(s.cfg.RequestSubscription, s)
}

func (s *Service) MessageDeserialized(deserialized []byte) {
	var request Message
	if err := json.Unmarshal(deserialized, &request); err != nil {
		logging.Warn(err, "unmarshal rag message")
		return
	}

	log.Info().Msgf("received message: %v", request)

	ctx, cancel := context.WithTimeout(s.cfg.Context, s.cfg.QueryTimeout)
	defer cancel()

	answer, err := s.query(ctx, request.Message)
	if err != nil {
		logging.Warn(err, "rag query")
		answer = fmt.Sprintf("query error: %v", err)
	}

	response := Message{
		IPAddress: request.IPAddress,
		Message:   answer,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		logging.Warn(err, "marshal rag response")
		return
	}

	log.Info().Msgf("publishing response: %v to topic: %s", response, s.cfg.ResponseTopic)

	logging.Warn(s.cfg.PubSub.Publish(s.cfg.ResponseTopic, payload), "publish rag response")
}

func (s *Service) MessageParseError(err error) {
	logging.Warn(err, "parse rag message")
}
