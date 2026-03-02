package messaging

import (
	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

type Serializer interface {
	Serialize(data any) ([]byte, error)
	Deserialize(data []byte, target any) error
}

type MessageProcessor struct {
	serializer Serializer
	compressor compression.Compressor
	encryptor  encryption.Encryptor

	enableSerialization bool
	enableCompression   bool
	enableEncryption    bool
}

func NewMessageProcessor(
	serializer Serializer,
	compressor compression.Compressor,
	encryptor encryption.Encryptor,
	enableSerialization, enableCompression, enableEncryption bool,
) *MessageProcessor {
	return &MessageProcessor{
		serializer:          serializer,
		compressor:          compressor,
		encryptor:           encryptor,
		enableSerialization: enableSerialization,
		enableCompression:   enableCompression,
		enableEncryption:    enableEncryption,
	}
}

func (p *MessageProcessor) Process(data []byte) ([]byte, error) {
	result := data
	var err error


	if p.enableCompression && p.compressor != nil {
		result, err = p.compressor.Compress(result)
		if err != nil {
			return nil, err
		}
	}

	if p.enableEncryption && p.encryptor != nil {
		result, err = p.encryptor.Encrypt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *MessageProcessor) Unprocess(data []byte) ([]byte, error) {
	result := data
	var err error

	if p.enableEncryption && p.encryptor != nil {
		result, err = p.encryptor.Decrypt(result)
		if err != nil {
			return nil, err
		}
	}

	if p.enableCompression && p.compressor != nil {
		result, err = p.compressor.Decompress(result)
		if err != nil {
			return nil, err
		}
	}


	return result, nil
}
