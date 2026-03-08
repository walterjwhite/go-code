package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/walterjwhite/go-code/lib/events"
	"github.com/walterjwhite/go-code/lib/events/transport"
	"github.com/walterjwhite/go-code/lib/io/serialization"
)

type MockPublisher struct {
	LastTopic string
	LastEvent any
	Err       error
	CloseErr  error
}

func (m *MockPublisher) Publish(ctx context.Context, topic string, event any) error {
	m.LastTopic = topic
	m.LastEvent = event
	return m.Err
}

func (m *MockPublisher) Close() error {
	return m.CloseErr
}

type MockSubscriber struct {
	Err          error
	CloseErr     error
	MessageBytes []byte
}

func (m *MockSubscriber) Subscribe(ctx context.Context, subscription string, handler transport.MessageHandler) error {
	if m.Err != nil {
		return m.Err
	}
	if m.MessageBytes != nil {
		return handler(ctx, m.MessageBytes)
	}
	return nil
}

func (m *MockSubscriber) Close() error {
	return m.CloseErr
}

func TestEventRegistry(t *testing.T) {
	registry := NewEventRegistry()

	t.Run("Register and Get Event", func(t *testing.T) {
		event := &events.Event{
			EventID: 1,
			Details: "Test event",
			SupportedActions: []events.Action{
				{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			},
		}

		err := registry.Register(event)
		require.NoError(t, err)

		retrieved, err := registry.Get(1)
		require.NoError(t, err)
		assert.Equal(t, event, retrieved)
	})

	t.Run("Duplicate Event Registration", func(t *testing.T) {
		registry := NewEventRegistry()
		event := &events.Event{EventID: 1, Details: "Test"}

		err := registry.Register(event)
		require.NoError(t, err)

		err = registry.Register(event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already registered")
	})

	t.Run("Get Non-existent Event", func(t *testing.T) {
		registry := NewEventRegistry()
		_, err := registry.Get(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("List Events", func(t *testing.T) {
		registry := NewEventRegistry()

		event1 := &events.Event{EventID: 1, Details: "Event 1"}
		event2 := &events.Event{EventID: 2, Details: "Event 2"}

		err := registry.Register(event1)
		require.NoError(t, err)
		err = registry.Register(event2)
		require.NoError(t, err)

		events := registry.List()
		assert.Equal(t, 2, len(events))
	})

	t.Run("Register Nil Event", func(t *testing.T) {
		registry := NewEventRegistry()
		err := registry.Register(nil)
		assert.Error(t, err)
	})

	t.Run("Register Event with Invalid ID", func(t *testing.T) {
		registry := NewEventRegistry()
		event := &events.Event{EventID: 0, Details: "Invalid"}
		err := registry.Register(event)
		assert.Error(t, err)
	})
}

func TestResponseHandler(t *testing.T) {
	registry := NewEventRegistry()
	mockPublisher := &MockPublisher{}
	serializer := serialization.NewJSONSerializer()
	handler := NewResponseHandler(registry, mockPublisher, serializer)

	event := &events.Event{
		EventID: 1,
		Details: "Test event",
		SupportedActions: []events.Action{
			{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			{ActionID: 2, Message: "Action 2", SupportsArgs: true},
		},
	}
	err := registry.Register(event)
	require.NoError(t, err)

	t.Run("Validate Valid Response", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 1,
		}
		err := handler.ValidateResponse(response)
		assert.NoError(t, err)
	})

	t.Run("Validate Response with Args", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 2,
			Args:     []string{"arg1", "arg2"},
		}
		err := handler.ValidateResponse(response)
		assert.NoError(t, err)
	})

	t.Run("Validate Response with Args on Non-supporting Action", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 1,
			Args:     []string{"arg1"},
		}
		err := handler.ValidateResponse(response)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not support arguments")
	})

	t.Run("Validate Invalid Action", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 999,
		}
		err := handler.ValidateResponse(response)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Validate Non-existent Event", func(t *testing.T) {
		response := &events.Response{
			EventID:  999,
			ActionID: 1,
		}
		err := handler.ValidateResponse(response)
		assert.Error(t, err)
	})

	t.Run("Publish Valid Response", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 1,
		}
		err := handler.PublishResponse(context.Background(), "responses", response)
		assert.NoError(t, err)
		assert.Equal(t, "responses", mockPublisher.LastTopic)
	})

	t.Run("Publish Invalid Response", func(t *testing.T) {
		response := &events.Response{
			EventID:  1,
			ActionID: 999,
		}
		err := handler.PublishResponse(context.Background(), "responses", response)
		assert.Error(t, err)
	})
}

func TestEventService(t *testing.T) {
	mockPublisher := &MockPublisher{}
	mockSubscriber := &MockSubscriber{}
	serializer := serialization.NewJSONSerializer()
	svc := NewEventService(mockPublisher, mockSubscriber, serializer)

	t.Run("Register Event", func(t *testing.T) {
		event := &events.Event{
			EventID: 1,
			Details: "Test event",
			SupportedActions: []events.Action{
				{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			},
		}

		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		retrieved, err := svc.GetEvent(1)
		require.NoError(t, err)
		assert.Equal(t, event, retrieved)
	})

	t.Run("List Events", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event1 := &events.Event{EventID: 1, Details: "Event 1"}
		event2 := &events.Event{EventID: 2, Details: "Event 2"}

		err := svc.RegisterEvent(event1)
		require.NoError(t, err)
		err = svc.RegisterEvent(event2)
		require.NoError(t, err)

		events := svc.ListEvents()
		assert.Equal(t, 2, len(events))
	})

	t.Run("Publish Event", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			Details: "Test event",
			SupportedActions: []events.Action{
				{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			},
		}

		err := svc.PublishEvent(context.Background(), "events", event)
		require.NoError(t, err)
		assert.Equal(t, "events", mockPublisher.LastTopic)
	})

	t.Run("Publish Response", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			Details: "Test event",
			SupportedActions: []events.Action{
				{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			},
		}
		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		response := &events.Response{
			EventID:  1,
			ActionID: 1,
		}

		err = svc.PublishResponse(context.Background(), "responses", response)
		require.NoError(t, err)
		assert.Equal(t, "responses", mockPublisher.LastTopic)
	})

	t.Run("Validate Response", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			Details: "Test event",
			SupportedActions: []events.Action{
				{ActionID: 1, Message: "Action 1", SupportsArgs: false},
			},
		}
		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		validResponse := &events.Response{
			EventID:  1,
			ActionID: 1,
		}
		err = svc.ValidateResponse(validResponse)
		assert.NoError(t, err)

		invalidResponse := &events.Response{
			EventID:  1,
			ActionID: 999,
		}
		err = svc.ValidateResponse(invalidResponse)
		assert.Error(t, err)
	})

	t.Run("Close Service", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		err := svc.Close()
		assert.NoError(t, err)
	})

	t.Run("Publish Event with Nil Event", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		err := svc.PublishEvent(context.Background(), "events", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event cannot be nil")
	})

	t.Run("Publish Event with Empty Topic", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{EventID: 1, Details: "Test"}
		err := svc.PublishEvent(context.Background(), "", event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "topic cannot be empty")
	})

	t.Run("Publish Response with Empty Topic", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			SupportedActions: []events.Action{
				{ActionID: 1, SupportsArgs: false},
			},
		}
		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		response := &events.Response{EventID: 1, ActionID: 1}
		err = svc.PublishResponse(context.Background(), "", response)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "topic cannot be empty")
	})

	t.Run("Publish Response with Long Topic", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			SupportedActions: []events.Action{
				{ActionID: 1, SupportsArgs: false},
			},
		}
		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		longTopic := string(make([]byte, 257))
		response := &events.Response{EventID: 1, ActionID: 1}
		err = svc.PublishResponse(context.Background(), longTopic, response)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "topic length exceeds maximum")
	})

	t.Run("Validate Nil Response", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		err := svc.ValidateResponse(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "response cannot be nil")
	})

	t.Run("Subscribe with Valid Handler", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		handler := func(event *events.Event) error {
			return nil
		}
		err := svc.Subscribe(context.Background(), "test-subscription", handler)
		assert.NoError(t, err)
	})

	t.Run("Subscribe with Nil Subscriber", func(t *testing.T) {
		svc := NewEventService(mockPublisher, nil, serializer)
		handler := func(event *events.Event) error {
			return nil
		}
		err := svc.Subscribe(context.Background(), "test-subscription", handler)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "subscriber not configured")
	})

	t.Run("Close Service with Errors", func(t *testing.T) {
		mockPub := &MockPublisher{CloseErr: assert.AnError}
		mockSub := &MockSubscriber{CloseErr: assert.AnError}
		svc := NewEventService(mockPub, mockSub, serializer)
		err := svc.Close()
		assert.Error(t, err)
	})

	t.Run("Publish Event Already Registered", func(t *testing.T) {
		svc := NewEventService(mockPublisher, mockSubscriber, serializer)
		event := &events.Event{EventID: 1, Details: "Test"}

		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		err = svc.PublishEvent(context.Background(), "events", event)
		assert.NoError(t, err)
	})

	t.Run("Publish Response with Publisher Error", func(t *testing.T) {
		mockPub := &MockPublisher{Err: assert.AnError}
		svc := NewEventService(mockPub, mockSubscriber, serializer)
		event := &events.Event{
			EventID: 1,
			SupportedActions: []events.Action{
				{ActionID: 1, SupportsArgs: false},
			},
		}
		err := svc.RegisterEvent(event)
		require.NoError(t, err)

		response := &events.Response{EventID: 1, ActionID: 1}
		err = svc.PublishResponse(context.Background(), "responses", response)
		assert.Error(t, err)
	})

	t.Run("Subscribe with Message Handling", func(t *testing.T) {
		event := &events.Event{EventID: 1, Details: "Test"}
		eventBytes, err := serializer.Serialize(event)
		require.NoError(t, err)

		mockSub := &MockSubscriber{MessageBytes: eventBytes}
		svc := NewEventService(mockPublisher, mockSub, serializer)

		handlerCalled := false
		handler := func(e *events.Event) error {
			handlerCalled = true
			assert.Equal(t, 1, e.EventID)
			return nil
		}

		err = svc.Subscribe(context.Background(), "test-subscription", handler)
		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})
}
