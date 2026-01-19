package google

import (
	"encoding/json"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walterjwhite/go-code/lib/io/compression/zstd"
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










