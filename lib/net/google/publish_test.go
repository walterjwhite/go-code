package google

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
)

func TestConf_serialize(t *testing.T) {
	message := []byte("test message")

	conf := &Conf{Serialize: false}
	serialized, err := conf.serialize(message)
	assert.NoError(t, err)
	assert.Equal(t, message, serialized)

	conf.Serialize = true
	serialized, err = conf.serialize(message)
	assert.NoError(t, err)
	expected, _ := json.Marshal(message)
	assert.Equal(t, expected, serialized)
}

func TestConf_encrypt(t *testing.T) {
	data := []byte("test data")

	conf := &Conf{aes: nil}
	encrypted := conf.encrypt(data)
	assert.Equal(t, data, encrypted)
}

func TestConf_compress(t *testing.T) {
	data := []byte("test data for compression")

	conf := &Conf{Compress: false}
	compressed := conf.compress(data)
	assert.Equal(t, data, compressed)

	conf.Compress = true
	compressed = conf.compress(data)
	assert.NotEqual(t, data, compressed) // Should be different after compression
	decompressed, err := zstd.DecompressBuffer(compressed)
	assert.NoError(t, err)
	assert.Equal(t, data, decompressed)
}

func TestConf_Publish_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	conf := &Conf{
		ctx:    ctx,
		client: nil, // Won't be used since context is cancelled
	}

	message := []byte("hello, pubsub!")
	err := conf.Publish("projects/test-project/topics/test-topic", message)
	assert.Error(t, err)
}

func TestConf_encrypt_WithAES(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	aesInstance, err := aes.New(key)
	assert.NoError(t, err)

	conf := &Conf{aes: aesInstance}
	data := []byte("test data to encrypt")

	encrypted := conf.encrypt(data)
	assert.NotEqual(t, data, encrypted, "encrypted data should differ from original")
	assert.NotEmpty(t, encrypted)
}

func TestConf_decrypt_WithAES(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	aesInstance, err := aes.New(key)
	assert.NoError(t, err)

	conf := &Conf{aes: aesInstance}
	data := []byte("test data to encrypt")

	encrypted := conf.encrypt(data)

	decrypted, err := conf.decrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, data, decrypted, "decrypted data should match original")
}

func TestConf_prepareMessage(t *testing.T) {
	conf := &Conf{
		Serialize: false,
		Compress:  false,
		aes:       nil,
	}

	message := []byte("test message")
	prepared, err := conf.prepareMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, message, prepared)

	conf.Serialize = true
	prepared, err = conf.prepareMessage(message)
	assert.NoError(t, err)
	assert.NotEqual(t, message, prepared)

	conf.Serialize = false
	conf.Compress = true
	prepared, err = conf.prepareMessage(message)
	assert.NoError(t, err)
	assert.NotEqual(t, message, prepared)

	decompressed, err := zstd.DecompressBuffer(prepared)
	assert.NoError(t, err)
	assert.Equal(t, message, decompressed)

	conf.Serialize = true
	conf.Compress = true
	prepared, err = conf.prepareMessage(message)
	assert.NoError(t, err)
	assert.NotEqual(t, message, prepared)
}
