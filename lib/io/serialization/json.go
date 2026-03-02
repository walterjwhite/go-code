package serialization

import (
	"encoding/json"
)

type JSONSerializer struct{}

var _ Serializer = (*JSONSerializer)(nil)

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (s *JSONSerializer) Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

func (s *JSONSerializer) Deserialize(data []byte, target any) error {
	return json.Unmarshal(data, target)
}
