package google

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option" // Re-added import

	"cloud.google.com/go/pubsub/pstest"
)

type MockMessageSubscriber struct {
	receivedMessages [][]byte
	parseErrors      []error
	mu               sync.Mutex
}

func (m *MockMessageSubscriber) MessageDeserialized(deserialized []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.receivedMessages = append(m.receivedMessages, deserialized)
}

func (m *MockMessageSubscriber) MessageParseError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.parseErrors = append(m.parseErrors, err)
}

func setupMockPubSubForSubscribe(ctx context.Context, t *testing.T) (*pubsub.Client, *pstest.Server, func()) {
	srv := pstest.NewServer()
	client, err := pubsub.NewClient(ctx, "project-id", option.WithoutAuthentication(), option.WithEndpoint(srv.Addr)) // Added options
	assert.NoError(t, err)

	return client, srv, func() {
		err := client.Close()
		if err != nil {
			t.Logf("Failed to close client: %v", err)
		}

		err = srv.Close()
		if err != nil {
			t.Logf("Failed to close server: %v", err)
		}
	}
}

func TestConf_decrypt(t *testing.T) {
	data := []byte("encrypted data")

	conf := &Conf{aes: nil}
	decrypted, err := conf.decrypt(data)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)

}

func TestConf_decompress(t *testing.T) {
	originalData := []byte("some data to decompress")

	conf := &Conf{Compress: false}
	decompressed := conf.decompress(originalData)
	assert.Equal(t, originalData, decompressed)

	conf.Compress = true
	compressedData := conf.compress(originalData) // Use the compress method from publish.go (same package)
	decompressed = conf.decompress(compressedData)
	assert.Equal(t, originalData, decompressed)

	conf.Compress = true
	invalidCompressedData := []byte{0x01, 0x02, 0x03}
	decompressed = conf.decompress(invalidCompressedData)
	assert.NotEqual(t, originalData, decompressed) // Should not equal original
	assert.Empty(t, decompressed)                  // Should return empty or error-like data if decompression fails
}

func TestConf_deserialize(t *testing.T) {
	mockSubscriber := &MockMessageSubscriber{}
	originalMessage := []byte("test message")

	conf := &Conf{Serialize: false}
	deserialized, err := conf.deserialize(mockSubscriber, originalMessage)
	assert.NoError(t, err)
	assert.Equal(t, originalMessage, deserialized)

	conf.Serialize = true
	jsonMarshaledMessage, _ := json.Marshal(originalMessage)
	deserialized, err = conf.deserialize(mockSubscriber, jsonMarshaledMessage)
	assert.NoError(t, err)
	assert.Equal(t, originalMessage, deserialized)
	assert.Len(t, mockSubscriber.receivedMessages, 0) // Should not be called directly by deserialize helper

	invalidJson := []byte(`{"key": "value"`)
	deserialized, err = conf.deserialize(mockSubscriber, invalidJson)
	assert.Error(t, err)
	assert.Nil(t, deserialized)
	assert.Len(t, mockSubscriber.parseErrors, 1) // Error should be logged by subscriber
}

func TestConf_Subscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockClient, srv, teardown := setupMockPubSubForSubscribe(ctx, t)
	defer teardown()
	_ = srv // Silence "declared and not used" warning for srv

	conf := &Conf{
		ctx:    ctx,
		client: mockClient,
	}

	mockSubscriber := &MockMessageSubscriber{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conf.Subscribe("test-topic-sub", "test-subscription", mockSubscriber)
	}()

	time.Sleep(100 * time.Millisecond)

	message := []byte("test message from publisher")
	msgID := srv.Publish("test-topic-sub", message, nil) // Use srv.Publish directly
	assert.NotEmpty(t, msgID)                            // Check if message ID is returned


	time.Sleep(500 * time.Millisecond)

	mockSubscriber.mu.Lock()
	assert.Len(t, mockSubscriber.receivedMessages, 1)
	assert.Equal(t, message, mockSubscriber.receivedMessages[0])
	mockSubscriber.mu.Unlock()

	conf.Compress = true
	conf.Serialize = true

	messageToPublish := []byte("secure message")
	var processedMessage = messageToPublish
	if conf.Serialize {
		processedMessage, _ = json.Marshal(processedMessage)
	}
	if conf.Compress {
		processedMessage = conf.compress(processedMessage)
	}
	msgID = srv.Publish("test-topic-sub", processedMessage, nil)
	assert.NotEmpty(t, msgID)

	time.Sleep(500 * time.Millisecond) // Wait for message to be processed

	mockSubscriber.mu.Lock()
	assert.Len(t, mockSubscriber.receivedMessages, 2) // One from before, one from now
	assert.Equal(t, messageToPublish, mockSubscriber.receivedMessages[1])
	mockSubscriber.mu.Unlock()

	cancel()
	wg.Wait() // Wait for the goroutine to finish
	fmt.Println("Subscriber goroutine finished.")
}
