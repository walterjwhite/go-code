package serialization

import (
	"testing"
)

func TestJSONSerializer(t *testing.T) {
	serializer := NewJSONSerializer()

	t.Run("Serialize and Deserialize", func(t *testing.T) {
		type TestData struct {
			Name  string
			Value int
		}

		original := TestData{Name: "test", Value: 42}

		data, err := serializer.Serialize(original)
		if err != nil {
			t.Fatalf("Serialize failed: %v", err)
		}

		var result TestData
		err = serializer.Deserialize(data, &result)
		if err != nil {
			t.Fatalf("Deserialize failed: %v", err)
		}

		if result.Name != original.Name || result.Value != original.Value {
			t.Errorf("Expected %+v, got %+v", original, result)
		}
	})

	t.Run("Serialize nil", func(t *testing.T) {
		data, err := serializer.Serialize(nil)
		if err != nil {
			t.Fatalf("Serialize nil failed: %v", err)
		}

		if string(data) != "null" {
			t.Errorf("Expected 'null', got %s", string(data))
		}
	})

	t.Run("Deserialize invalid JSON", func(t *testing.T) {
		var result map[string]any
		err := serializer.Deserialize([]byte("invalid json"), &result)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}
