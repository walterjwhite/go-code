package messaging

import (
	"errors"
	"fmt"

	"github.com/walterjwhite/go-code/lib/io/compression"
	"github.com/walterjwhite/go-code/lib/security/encryption"
)

const DefaultMaxMessageSize = 10 * 1024 * 1024

type Serializer interface {
	Serialize(data any) ([]byte, error)
	Deserialize(data []byte, target any) error
}

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

type MessageProcessorConfig struct {
	MaxMessageSize int
	EnableLogging  bool
}

func DefaultConfig() MessageProcessorConfig {
	return MessageProcessorConfig{
		MaxMessageSize: DefaultMaxMessageSize,
		EnableLogging:  false,
	}
}

type MessageProcessor struct {
	serializer Serializer
	compressor compression.Compressor
	encryptor  encryption.Encryptor
	logger     Logger

	enableSerialization bool
	enableCompression   bool
	enableEncryption    bool
	maxMessageSize      int
}

func NewMessageProcessor(
	serializer Serializer,
	compressor compression.Compressor,
	encryptor encryption.Encryptor,
	enableSerialization, enableCompression, enableEncryption bool,
) *MessageProcessor {
	return NewMessageProcessorWithConfig(
		serializer,
		compressor,
		encryptor,
		enableSerialization,
		enableCompression,
		enableEncryption,
		DefaultConfig(),
		nil,
	)
}

func NewMessageProcessorWithConfig(
	serializer Serializer,
	compressor compression.Compressor,
	encryptor encryption.Encryptor,
	enableSerialization, enableCompression, enableEncryption bool,
	config MessageProcessorConfig,
	logger Logger,
) *MessageProcessor {
	effectiveEncryption := enableEncryption

	maxSize := config.MaxMessageSize
	if maxSize <= 0 {
		maxSize = DefaultMaxMessageSize
	}

	return &MessageProcessor{
		serializer:          serializer,
		compressor:          compressor,
		encryptor:           encryptor,
		logger:              logger,
		enableSerialization: enableSerialization,
		enableCompression:   enableCompression,
		enableEncryption:    effectiveEncryption,
		maxMessageSize:      maxSize,
	}
}

func (p *MessageProcessor) validateInput(data []byte, operation string) error {
	if data == nil {
		if p.logger != nil {
			p.logger.Error("nil input data", "operation", operation)
		}
		return errors.New("input data cannot be nil")
	}

	if len(data) == 0 {
		if p.logger != nil {
			p.logger.Error("empty input data", "operation", operation)
		}
		return errors.New("input data cannot be empty")
	}

	if len(data) > p.maxMessageSize {
		if p.logger != nil {
			p.logger.Error("input data exceeds maximum size",
				"operation", operation,
				"size", len(data),
				"maxSize", p.maxMessageSize)
		}
		return fmt.Errorf("input data size %d exceeds maximum allowed size %d", len(data), p.maxMessageSize)
	}

	return nil
}

func (p *MessageProcessor) Process(data []byte) ([]byte, error) {
	if err := p.validateInput(data, "Process"); err != nil {
		return nil, err
	}

	result := data
	var err error

	if p.logger != nil {
		p.logger.Info("processing message",
			"initialSize", len(data),
			"compressionEnabled", p.enableCompression,
			"encryptionEnabled", p.enableEncryption)
	}

	if p.enableCompression && p.compressor != nil {
		result, err = p.compressor.Compress(result)
		if err != nil {
			if p.logger != nil {
				p.logger.Error("compression failed", "error", err.Error())
			}
			return nil, fmt.Errorf("compression failed: %w", err)
		}
	}

	if p.enableEncryption && p.encryptor != nil {
		result, err = p.encryptor.Encrypt(result)
		if err != nil {
			if p.logger != nil {
				p.logger.Error("encryption failed", "error", err.Error())
			}
			return nil, fmt.Errorf("encryption failed: %w", err)
		}
	}

	if p.logger != nil {
		p.logger.Info("message processed successfully",
			"finalSize", len(result),
			"compressionRatio", float64(len(result))/float64(len(data)))
	}

	return result, nil
}

func (p *MessageProcessor) Unprocess(data []byte) ([]byte, error) {
	if err := p.validateInput(data, "Unprocess"); err != nil {
		return nil, err
	}

	result := data
	var err error

	if p.logger != nil {
		p.logger.Info("unprocessing message",
			"initialSize", len(data),
			"encryptionEnabled", p.enableEncryption,
			"compressionEnabled", p.enableCompression)
	}

	if p.enableEncryption && p.encryptor != nil {
		result, err = p.encryptor.Decrypt(result)
		if err != nil {
			if p.logger != nil {
				p.logger.Error("decryption failed", "error", err.Error())
			}
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
	}

	if p.enableCompression && p.compressor != nil {
		result, err = p.compressor.Decompress(result)
		if err != nil {
			if p.logger != nil {
				p.logger.Error("decompression failed", "error", err.Error())
			}
			return nil, fmt.Errorf("decompression failed: %w", err)
		}
	}

	if p.logger != nil {
		p.logger.Info("message unprocessed successfully",
			"finalSize", len(result))
	}

	return result, nil
}
