package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/walterjwhite/go-code/lib/events"
	"github.com/walterjwhite/go-code/lib/events/transport"
	"github.com/walterjwhite/go-code/lib/io/serialization"
)

type EventRegistry struct {
	events map[int]*events.Event
	mu     sync.RWMutex
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		events: make(map[int]*events.Event),
	}
}

func (r *EventRegistry) Register(event *events.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if event == nil {
		return fmt.Errorf("cannot register nil event")
	}

	if event.EventID <= 0 {
		return fmt.Errorf("event ID must be positive")
	}

	if _, exists := r.events[event.EventID]; exists {
		return fmt.Errorf("event ID %d is already registered", event.EventID)
	}

	r.events[event.EventID] = event
	return nil
}

func (r *EventRegistry) Get(eventID int) (*events.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[eventID]
	if !exists {
		return nil, fmt.Errorf("event ID %d not found", eventID)
	}

	return event, nil
}

func (r *EventRegistry) List() []*events.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*events.Event, 0, len(r.events))
	for _, event := range r.events {
		eventCopy := *event
		result = append(result, &eventCopy)
	}
	return result
}

type ResponseHandler struct {
	registry   *EventRegistry
	publisher  transport.Publisher
	serializer serialization.Serializer
}

func NewResponseHandler(
	registry *EventRegistry,
	publisher transport.Publisher,
	serializer serialization.Serializer,
) *ResponseHandler {
	return &ResponseHandler{
		registry:   registry,
		publisher:  publisher,
		serializer: serializer,
	}
}

func (h *ResponseHandler) ValidateResponse(response *events.Response) error {
	if response == nil {
		return fmt.Errorf("response cannot be nil")
	}

	event, err := h.registry.Get(response.EventID)
	if err != nil {
		return err
	}

	var action *events.Action
	for _, a := range event.SupportedActions {
		if a.ActionID == response.ActionID {
			action = &a
			break
		}
	}

	if action == nil {
		return fmt.Errorf("action ID %d not found for event %d", response.ActionID, response.EventID)
	}

	if len(response.Args) > 0 && !action.SupportsArgs {
		return fmt.Errorf("action %d does not support arguments", response.ActionID)
	}

	return nil
}

func (h *ResponseHandler) PublishResponse(ctx context.Context, topic string, response *events.Response) error {
	if topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}
	if len(topic) > 256 {
		return fmt.Errorf("topic length exceeds maximum allowed (256 chars)")
	}

	if err := h.ValidateResponse(response); err != nil {
		return fmt.Errorf("invalid response: %w", err)
	}

	if err := h.publisher.Publish(ctx, topic, response); err != nil {
		return fmt.Errorf("failed to publish response: %w", err)
	}

	return nil
}

type EventService struct {
	registry   *EventRegistry
	handler    *ResponseHandler
	publisher  transport.Publisher
	subscriber transport.Subscriber
	serializer serialization.Serializer
}

func NewEventService(
	publisher transport.Publisher,
	subscriber transport.Subscriber,
	serializer serialization.Serializer,
) *EventService {
	registry := NewEventRegistry()
	handler := NewResponseHandler(registry, publisher, serializer)

	return &EventService{
		registry:   registry,
		handler:    handler,
		publisher:  publisher,
		subscriber: subscriber,
		serializer: serializer,
	}
}

func (s *EventService) RegisterEvent(event *events.Event) error {
	return s.registry.Register(event)
}

func (s *EventService) GetEvent(eventID int) (*events.Event, error) {
	return s.registry.Get(eventID)
}

func (s *EventService) ListEvents() []*events.Event {
	return s.registry.List()
}

func (s *EventService) PublishEvent(ctx context.Context, topic string, event *events.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}
	if topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}

	err := s.registry.Register(event)
	if err != nil {
		_, getErr := s.registry.Get(event.EventID)
		if getErr != nil {
			return err
		}
	}

	return s.publisher.Publish(ctx, topic, event)
}

func (s *EventService) PublishResponse(ctx context.Context, topic string, response *events.Response) error {
	return s.handler.PublishResponse(ctx, topic, response)
}

func (s *EventService) ValidateResponse(response *events.Response) error {
	return s.handler.ValidateResponse(response)
}

func (s *EventService) Subscribe(ctx context.Context, subscription string, handler func(event *events.Event) error) error {
	if s.subscriber == nil {
		return fmt.Errorf("subscriber not configured")
	}

	messageHandler := func(ctx context.Context, message []byte) error {
		var event events.Event
		if err := s.serializer.Deserialize(message, &event); err != nil {
			return fmt.Errorf("failed to deserialize event: %w", err)
		}

		return handler(&event)
	}

	return s.subscriber.Subscribe(ctx, subscription, messageHandler)
}

func (s *EventService) Close() error {
	var errors []error

	if s.publisher != nil {
		if err := s.publisher.Close(); err != nil {
			errors = append(errors, fmt.Errorf("publisher close error: %w", err))
		}
	}

	if s.subscriber != nil {
		if err := s.subscriber.Close(); err != nil {
			errors = append(errors, fmt.Errorf("subscriber close error: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("service close errors: %v", errors)
	}

	return nil
}
