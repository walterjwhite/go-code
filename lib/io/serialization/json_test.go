package serialization

import (
	"context"
	"strings"
	"testing"
	"time"
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

func TestJSONSerializerSecurity(t *testing.T) {
	t.Run("Deserialize empty data", func(t *testing.T) {
		serializer := NewJSONSerializer()
		var result map[string]any
		err := serializer.Deserialize([]byte{}, &result)
		if err == nil {
			t.Error("Expected error for empty data, got nil")
		}
		if !strings.Contains(err.Error(), "empty input data") {
			t.Errorf("Expected 'empty input data' error, got: %v", err)
		}
	})

	t.Run("Deserialize exceeds max size", func(t *testing.T) {
		serializer := NewJSONSerializerWithMaxSize(10)
		largeData := []byte(`{"data":"` + strings.Repeat("a", 100) + `"}`)
		var result map[string]any
		err := serializer.Deserialize(largeData, &result)
		if err == nil {
			t.Error("Expected error for data exceeding max size, got nil")
		}
		if !strings.Contains(err.Error(), "exceeds maximum allowed limit") {
			t.Errorf("Expected size exceeded error, got: %v", err)
		}
	})

	t.Run("Deserialize with nil target", func(t *testing.T) {
		serializer := NewJSONSerializer()
		err := serializer.Deserialize([]byte(`{"key":"value"}`), nil)
		if err == nil {
			t.Error("Expected error for nil target, got nil")
		}
		if !strings.Contains(err.Error(), "target cannot be nil") {
			t.Errorf("Expected 'target cannot be nil' error, got: %v", err)
		}
	})

	t.Run("SerializeWithContext with timeout", func(t *testing.T) {
		serializer := NewJSONSerializer()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		type TestData struct {
			Name string
		}
		data := TestData{Name: "test"}

		result, err := serializer.SerializeWithContext(ctx, data)
		if err != nil {
			t.Fatalf("SerializeWithContext failed: %v", err)
		}
		if string(result) != `{"Name":"test"}` {
			t.Errorf("Unexpected result: %s", string(result))
		}
	})

	t.Run("DeserializeWithContext with timeout", func(t *testing.T) {
		serializer := NewJSONSerializer()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		var result map[string]any
		err := serializer.DeserializeWithContext(ctx, []byte(`{"key":"value"}`), &result)
		if err != nil {
			t.Fatalf("DeserializeWithContext failed: %v", err)
		}
		if result["key"] != "value" {
			t.Errorf("Expected 'value', got %v", result["key"])
		}
	})

	t.Run("DeserializeWithContext with cancelled context", func(t *testing.T) {
		serializer := NewJSONSerializer()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var result map[string]any
		err := serializer.DeserializeWithContext(ctx, []byte(`{"key":"value"}`), &result)
		if err == nil {
			t.Error("Expected error for cancelled context, got nil")
		}
		if !strings.Contains(err.Error(), "cancelled") {
			t.Errorf("Expected 'cancelled' error, got: %v", err)
		}
	})

	t.Run("Custom max size configuration", func(t *testing.T) {
		customMaxSize := int64(1024)
		serializer := NewJSONSerializerWithMaxSize(customMaxSize)

		validData := []byte(`{"small":"data"}`)
		var result map[string]any
		err := serializer.Deserialize(validData, &result)
		if err != nil {
			t.Fatalf("Expected success for data within limit, got: %v", err)
		}

		invalidData := []byte(`{"` + strings.Repeat("a", 2000) + `":"value"}`)
		err = serializer.Deserialize(invalidData, &result)
		if err == nil {
			t.Error("Expected error for data exceeding custom limit, got nil")
		}
	})
}
