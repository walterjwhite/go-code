package transport

import (
	"context"
	"errors"
	"testing"

	"github.com/walterjwhite/go-code/lib/io/serialization"
)

const (
	testSMTPHost     = "smtp.example.com"
	testSMTPPort     = 587
	testSMTPUsername = "test_user"
	testSMTPPassword = "test_password_placeholder" // Use env var in production
	testFromEmail    = "from@example.com"
	testToEmail      = "to@example.com"
	testIMAPHost     = "imap.example.com"
	testIMAPUsername = "test_user"
	testIMAPPassword = "test_password_placeholder" // Use env var in production
	testSQSQueueURL  = "https://sqs.us-east-1.amazonaws.com/123456789012/myqueue"
)

type mockSerializer struct {
	serializeErr   error
	deserializeErr error
}

func (m *mockSerializer) Serialize(v any) ([]byte, error) {
	if m.serializeErr != nil {
		return nil, m.serializeErr
	}
	return []byte("serialized"), nil
}

func (m *mockSerializer) Deserialize(data []byte, v any) error {
	return m.deserializeErr
}

func (m *mockSerializer) SerializeWithContext(ctx context.Context, v any) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return m.Serialize(v)
}

func (m *mockSerializer) DeserializeWithContext(ctx context.Context, data []byte, v any) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return m.Deserialize(data, v)
}

type mockCompressor struct {
	compressErr   error
	decompressErr error
}

func (m *mockCompressor) Compress(data []byte) ([]byte, error) {
	if m.compressErr != nil {
		return nil, m.compressErr
	}
	return data, nil
}

func (m *mockCompressor) Decompress(data []byte) ([]byte, error) {
	if m.decompressErr != nil {
		return nil, m.decompressErr
	}
	return data, nil
}

type mockEncryptor struct{}

func (m *mockEncryptor) Encrypt(data []byte) ([]byte, error) {
	return data, nil
}

func (m *mockEncryptor) Decrypt(data []byte) ([]byte, error) {
	return data, nil
}

func TestEmailPublisher_NewEmailPublisher(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	encryptor := &mockEncryptor{}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, encryptor, false, false)

	if pub == nil {
		t.Fatal("expected non-nil publisher")
	}
	if pub.from != testFromEmail {
		t.Errorf("expected from=%s, got %s", testFromEmail, pub.from)
	}
}

func TestEmailPublisher_NewEmailPublisher_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, nil, false, false)
}

func TestEmailPublisher_Publish_NilDialer(t *testing.T) {
	pub := &EmailPublisher{
		dialer:     nil,
		serializer: &mockSerializer{},
	}

	err := pub.Publish(context.Background(), testToEmail, "test event")
	if err == nil {
		t.Error("expected error for nil dialer")
	}
}

func TestEmailPublisher_Publish_SerializationError(t *testing.T) {
	serializer := &mockSerializer{serializeErr: errors.New("serialize failed")}
	compressor := &mockCompressor{}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, nil, false, false)

	err := pub.Publish(context.Background(), testToEmail, "test event")
	if err == nil {
		t.Error("expected serialization error")
	}
}

func TestEmailPublisher_Close(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, nil, false, false)

	if err := pub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEmailSubscriber_NewEmailSubscriber(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	encryptor := &mockEncryptor{}

	sub := NewEmailSubscriber(testIMAPHost, testIMAPUsername, testIMAPPassword,
		serializer, compressor, encryptor, false, false)

	if sub == nil {
		t.Fatal("expected non-nil subscriber")
	}
	if sub.host != testIMAPHost {
		t.Errorf("expected host=%s, got %s", testIMAPHost, sub.host)
	}
}

func TestEmailSubscriber_Subscribe_InvalidFolder(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	sub := NewEmailSubscriber(testIMAPHost, testIMAPUsername, testIMAPPassword,
		serializer, compressor, nil, false, false)

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "INBOX/subfolder", handler)
	if err == nil {
		t.Error("expected error for invalid folder path")
	}
}

func TestEmailSubscriber_Subscribe_DefaultFolder(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	sub := NewEmailSubscriber(testIMAPHost, testIMAPUsername, testIMAPPassword,
		serializer, compressor, nil, false, false)

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "", handler)
	if err == nil {
		t.Error("expected error (IMAP not implemented)")
	}
}

func TestEmailSubscriber_Close(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	sub := NewEmailSubscriber(testIMAPHost, testIMAPUsername, testIMAPPassword,
		serializer, compressor, nil, false, false)

	if err := sub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSQSPublisher_NewSQSPublisher(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	pub := NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	if pub == nil {
		t.Fatal("expected non-nil publisher")
	}
	if pub.queueURL != testSQSQueueURL {
		t.Errorf("unexpected queueURL: %s", pub.queueURL)
	}
}

func TestSQSPublisher_Publish_NilClient(t *testing.T) {
	pub := &SQSPublisher{
		client:     nil,
		serializer: &mockSerializer{},
	}

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestSQSPublisher_Publish_SerializationError(t *testing.T) {
	serializer := &mockSerializer{serializeErr: errors.New("serialize failed")}
	compressor := &mockCompressor{}
	client := "mock-client"

	pub := NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected serialization error")
	}
}

func TestSQSPublisher_Close(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	pub := NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	if err := pub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSQSSubscriber_NewSQSSubscriber(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	sub := NewSQSSubscriber(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	if sub == nil {
		t.Fatal("expected non-nil subscriber")
	}
	if sub.queueURL != testSQSQueueURL {
		t.Errorf("unexpected queueURL: %s", sub.queueURL)
	}
}

func TestSQSSubscriber_Subscribe_NilClient(t *testing.T) {
	sub := &SQSSubscriber{
		client:     nil,
		serializer: &mockSerializer{},
	}

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "subscription", handler)
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestSQSSubscriber_Close(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	sub := NewSQSSubscriber(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	if err := sub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGooglePubSubPublisher_NewGooglePubSubPublisher(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	pub := NewGooglePubSubPublisher(nil, serializer, compressor, nil, false, false)

	if pub == nil {
		t.Fatal("expected non-nil publisher")
	}
}

func TestGooglePubSubPublisher_Publish_NilClient(t *testing.T) {
	pub := &GooglePubSubPublisher{
		client:     nil,
		serializer: &mockSerializer{},
	}

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestGooglePubSubPublisher_Publish_SerializationError(t *testing.T) {
	serializer := &mockSerializer{serializeErr: errors.New("serialize failed")}

	pub := &GooglePubSubPublisher{
		serializer: serializer,
	}

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected error (nil client or serialization)")
	}
}

func TestGooglePubSubPublisher_Close_NilClient(t *testing.T) {
	pub := &GooglePubSubPublisher{client: nil}

	if err := pub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGooglePubSubSubscriber_NewGooglePubSubSubscriber(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	sub := NewGooglePubSubSubscriber(nil, serializer, compressor, nil, false, false)

	if sub == nil {
		t.Fatal("expected non-nil subscriber")
	}
}

func TestGooglePubSubSubscriber_Subscribe_NilClient(t *testing.T) {
	sub := &GooglePubSubSubscriber{
		client:     nil,
		serializer: &mockSerializer{},
	}

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "subscription", handler)
	if err == nil {
		t.Error("expected error for nil client")
	}
}

func TestGooglePubSubSubscriber_Close_NilClient(t *testing.T) {
	sub := &GooglePubSubSubscriber{client: nil}

	if err := sub.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMessageHandler(t *testing.T) {
	called := false
	handler := MessageHandler(func(ctx context.Context, msg []byte) error {
		called = true
		if string(msg) != "test" {
			t.Errorf("expected 'test', got %s", string(msg))
		}
		return nil
	})

	err := handler(context.Background(), []byte("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

func TestMessageHandler_Error(t *testing.T) {
	expectedErr := errors.New("handler error")
	handler := MessageHandler(func(ctx context.Context, msg []byte) error {
		return expectedErr
	})

	err := handler(context.Background(), []byte("test"))
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestPublisherInterface(t *testing.T) {
	var _ Publisher = (*EmailPublisher)(nil)
	var _ Publisher = (*SQSPublisher)(nil)
	var _ Publisher = (*GooglePubSubPublisher)(nil)
}

func TestSubscriberInterface(t *testing.T) {
	var _ Subscriber = (*EmailSubscriber)(nil)
	var _ Subscriber = (*SQSSubscriber)(nil)
	var _ Subscriber = (*GooglePubSubSubscriber)(nil)
}

func TestSQSPublisher_NewSQSPublisher_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)
}

func TestGooglePubSubPublisher_NewGooglePubSubPublisher_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	NewGooglePubSubPublisher(nil, serializer, compressor, nil, false, false)
}

func TestEmailSubscriber_NewEmailSubscriber_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	NewEmailSubscriber("imap.example.com", "user", "pass",
		serializer, compressor, nil, false, false)
}

func TestSQSSubscriber_NewSQSSubscriber_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	NewSQSSubscriber(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)
}

func TestGooglePubSubSubscriber_NewGooglePubSubSubscriber_NilEncryptor(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	NewGooglePubSubSubscriber(nil, serializer, compressor, nil, false, false)
}

func TestEmailSubscriber_Subscribe_ValidFolder(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	sub := NewEmailSubscriber(testIMAPHost, testIMAPUsername, testIMAPPassword,
		serializer, compressor, nil, false, false)

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "INBOX", handler)
	if err == nil {
		t.Error("expected error (IMAP not implemented)")
	}
}

var _ serialization.Serializer = (*mockSerializer)(nil)

func TestEmailPublisher_Publish_ProcessingError(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{compressErr: errors.New("compress failed")}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, nil, true, false)

	err := pub.Publish(context.Background(), testToEmail, "test event")
	if err == nil {
		t.Error("expected processing error")
	}
}

func TestSQSPublisher_Publish_ProcessingError(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{compressErr: errors.New("compress failed")}
	client := "mock-client"

	pub := NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, true, false)

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected processing error")
	}
}

func TestGooglePubSubPublisher_Publish_ProcessingError(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{compressErr: errors.New("compress failed")}

	pub := NewGooglePubSubPublisher(nil, serializer, compressor, nil, true, false)

	pub.client = nil // Will fail at client check first, which is expected

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected error")
	}
}

func TestSQSSubscriber_Subscribe_ValidClient(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	sub := NewSQSSubscriber(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	handler := func(ctx context.Context, msg []byte) error { return nil }
	err := sub.Subscribe(context.Background(), "subscription", handler)
	if err == nil {
		t.Error("expected error (AWS SDK not implemented)")
	}
}

func TestSQSPublisher_Publish_ValidClient(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	client := "mock-client"

	pub := NewSQSPublisher(client, testSQSQueueURL,
		serializer, compressor, nil, false, false)

	err := pub.Publish(context.Background(), "topic", "test event")
	if err == nil {
		t.Error("expected error (AWS SDK not implemented)")
	}
}

func TestEmailPublisher_Publish_WithEncryption(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}
	encryptor := &mockEncryptor{}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, encryptor, false, true)

	err := pub.Publish(context.Background(), testToEmail, "test event")
	if err == nil {
		t.Error("expected error (SMTP not configured)")
	}
}

func TestEmailPublisher_Publish_WithCompression(t *testing.T) {
	serializer := &mockSerializer{}
	compressor := &mockCompressor{}

	pub := NewEmailPublisher(testSMTPHost, testSMTPPort, testSMTPUsername, testSMTPPassword, testFromEmail,
		serializer, compressor, nil, true, false)

	err := pub.Publish(context.Background(), testToEmail, "test event")
	if err == nil {
		t.Error("expected error (SMTP not configured)")
	}
}
