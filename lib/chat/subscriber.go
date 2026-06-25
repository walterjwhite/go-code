package chat

import "fmt"

func (s *responseSubscriber) MessageDeserialized(deserialized []byte) {
	message, err := s.deserialize(deserialized)
	if err != nil {
		s.pushError(err)
		return
	}
	s.pushResponse(message)
}

func (s *responseSubscriber) deserialize(payload []byte) (ServerMessage, error) {
	var message ServerMessage
	err := s.serializer.Deserialize(payload, &message)
	if err != nil {
		return ServerMessage{}, fmt.Errorf("deserialize server message: %w", err)
	}
	return message, nil
}

func (s *responseSubscriber) pushResponse(message ServerMessage) {
	select {
	case <-s.ctx.Done():
		s.pushError(s.ctx.Err())
	case s.responses <- message:
	}
}

func (s *responseSubscriber) MessageParseError(err error) {
	s.pushError(fmt.Errorf("pubsub message parse error: %w", err))
}

func (s *responseSubscriber) pushError(err error) {
	select {
	case s.errors <- err:
	default:
	}
}
