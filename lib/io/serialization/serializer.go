package serialization

import "context"

const DefaultMaxJSONSize = 10 * 1024 * 1024

type Serializer interface {
	Serialize(data any) ([]byte, error)
	SerializeWithContext(ctx context.Context, data any) ([]byte, error)
	Deserialize(data []byte, target any) error
	DeserializeWithContext(ctx context.Context, data []byte, target any) error
}
