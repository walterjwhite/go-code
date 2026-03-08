package serialization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var ErrJSONSizeExceeded = errors.New("json payload size exceeds maximum allowed limit")

type JSONSerializer struct {
	maxSize int64
}

var _ Serializer = (*JSONSerializer)(nil)

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{
		maxSize: DefaultMaxJSONSize,
	}
}

func NewJSONSerializerWithMaxSize(maxSize int64) *JSONSerializer {
	return &JSONSerializer{
		maxSize: maxSize,
	}
}

func (s *JSONSerializer) Serialize(data any) ([]byte, error) {
	return s.SerializeWithContext(context.Background(), data)
}

func (s *JSONSerializer) SerializeWithContext(ctx context.Context, data any) ([]byte, error) {
	if data == nil {
		return []byte("null"), nil
	}

	resultChan := make(chan struct {
		data []byte
		err  error
	}, 1)

	go func() {
		result, err := json.Marshal(data)
		resultChan <- struct {
			data []byte
			err  error
		}{result, err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("serialization cancelled: %w", ctx.Err())
	case result := <-resultChan:
		return result.data, result.err
	}
}

func (s *JSONSerializer) Deserialize(data []byte, target any) error {
	return s.DeserializeWithContext(context.Background(), data, target)
}

func (s *JSONSerializer) DeserializeWithContext(ctx context.Context, data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty input data")
	}

	if int64(len(data)) > s.maxSize {
		return fmt.Errorf("%w: size %d exceeds limit %d", ErrJSONSizeExceeded, len(data), s.maxSize)
	}

	if target == nil {
		return errors.New("target cannot be nil")
	}

	errChan := make(chan error, 1)

	go func() {
		err := json.Unmarshal(data, target)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("deserialization cancelled: %w", ctx.Err())
	case err := <-errChan:
		return err
	}
}
