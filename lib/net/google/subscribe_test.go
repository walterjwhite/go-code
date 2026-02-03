package google

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
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


func TestConf_decrypt(t *testing.T) {
	data := []byte("encrypted data")

	conf := &Conf{aes: nil}
	decrypted, err := conf.decrypt(data)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)

	key := []byte("01234567890123456789012345678901")
	aesInstance, err := aes.New(key)
	assert.NoError(t, err)

	conf.aes = aesInstance
	originalData := []byte("test data to decrypt")
	encrypted := conf.encrypt(originalData)
	
	decrypted, err = conf.decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, originalData, decrypted)
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
	
	conf := &Conf{aes: nil}
	data := []byte("test data")
	decrypted, err := conf.decrypt(data)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted)

	conf.Compress = false
	decompressed := conf.decompress(data)
	assert.Equal(t, data, decompressed)

	conf.Compress = true
	originalData := []byte("some data to decompress")
	compressedData := conf.compress(originalData)
	decompressed = conf.decompress(compressedData)
	assert.Equal(t, originalData, decompressed)

	mockSubscriber := &MockMessageSubscriber{}
	conf.Serialize = false
	deserialized, err := conf.deserialize(mockSubscriber, data)
	assert.NoError(t, err)
	assert.Equal(t, data, deserialized)

	conf.Serialize = true
	jsonMarshaledMessage, _ := json.Marshal(originalData)
	deserialized, err = conf.deserialize(mockSubscriber, jsonMarshaledMessage)
	assert.NoError(t, err)
	assert.Equal(t, originalData, deserialized)
}

func TestConf_processMessage(t *testing.T) {
	mockSubscriber := &MockMessageSubscriber{}
	conf := &Conf{
		Serialize: false,
		Compress:  false,
		aes:       nil,
	}

	message := []byte("test message")
	err := conf.processMessage(mockSubscriber, message)
	assert.NoError(t, err)
	assert.Len(t, mockSubscriber.receivedMessages, 1)
	assert.Equal(t, message, mockSubscriber.receivedMessages[0])

	mockSubscriber2 := &MockMessageSubscriber{}
	conf.Serialize = true
	jsonMessage, _ := json.Marshal(message)
	err = conf.processMessage(mockSubscriber2, jsonMessage)
	assert.NoError(t, err)
	assert.Len(t, mockSubscriber2.receivedMessages, 1)
	assert.Equal(t, message, mockSubscriber2.receivedMessages[0])

	mockSubscriber3 := &MockMessageSubscriber{}
	conf.Serialize = false
	conf.Compress = true
	compressedMessage := conf.compress(message)
	err = conf.processMessage(mockSubscriber3, compressedMessage)
	assert.NoError(t, err)
	assert.Len(t, mockSubscriber3.receivedMessages, 1)
	assert.Equal(t, message, mockSubscriber3.receivedMessages[0])

	mockSubscriber4 := &MockMessageSubscriber{}
	conf.Serialize = true
	conf.Compress = true
	jsonMsg, _ := json.Marshal(message)
	compressedJsonMsg := conf.compress(jsonMsg)
	err = conf.processMessage(mockSubscriber4, compressedJsonMsg)
	assert.NoError(t, err)
	assert.Len(t, mockSubscriber4.receivedMessages, 1)
	assert.Equal(t, message, mockSubscriber4.receivedMessages[0])

	mockSubscriber5 := &MockMessageSubscriber{}
	conf.Serialize = true
	conf.Compress = false
	invalidJSON := []byte(`{"incomplete": json`)
	err = conf.processMessage(mockSubscriber5, invalidJSON)
	assert.Error(t, err)
	assert.Len(t, mockSubscriber5.parseErrors, 1)
}
