package interfaces

import (
	"context"
	"io"
	"net/url"
	"regexp"
	"strings"
)

type Serializer interface {
	Serialize(data any) ([]byte, error)
	Deserialize(data []byte, target any) error
}

type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type Encryptor interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type MessageProcessor interface {
	Process(data []byte) ([]byte, error)
	Unprocess(data []byte) ([]byte, error)
}

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Close() error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

type MessageHandler interface {
	Handle(ctx context.Context, message []byte) error
	HandleError(ctx context.Context, err error)
}

type EmailSender interface {
	Send(ctx context.Context, message EmailMessage) error
}

type EmailMessage interface {
	From() string
	To() []string
	Subject() string
	Body() string
	Attachments() []Attachment
}

type Attachment interface {
	Name() string
	Content() io.Reader
	ContentType() string
}

type ConfigLoader interface {
	Load(target any) error
	LoadFrom(source string, target any) error
}

type SecretProvider interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

type WorkScheduler interface {
	Schedule(ctx context.Context, work Work) error
	Cancel(workID string) error
	Status(workID string) (WorkStatus, error)
}

type Work interface {
	ID() string
	Execute(ctx context.Context) error
	OnComplete(ctx context.Context) error
	OnError(ctx context.Context, err error) error
}

type WorkStatus int

const (
	WorkStatusPending WorkStatus = iota
	WorkStatusRunning
	WorkStatusCompleted
	WorkStatusFailed
	WorkStatusCancelled
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl int64) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	Post(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Put(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error)
}

type MetricsCollector interface {
	Increment(metric string, tags map[string]string)
	Gauge(metric string, value float64, tags map[string]string)
	Histogram(metric string, value float64, tags map[string]string)
	Timing(metric string, duration int64, tags map[string]string)
}

type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

type SecureSecretProvider interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	Rotate(ctx context.Context, key string, newValue string) error
}

type SanitizedField struct {
	Key   string
	Value string
}

func Sanitize(value string) string {
	sanitized := strings.ReplaceAll(value, "\n", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	if len(sanitized) > 10000 {
		sanitized = sanitized[:10000] + "...[truncated]"
	}
	return sanitized
}

type SecureLogger interface {
	Debug(msg string, fields ...SanitizedField)
	Info(msg string, fields ...SanitizedField)
	Warn(msg string, fields ...SanitizedField)
	Error(msg string, fields ...SanitizedField)
	Fatal(msg string, fields ...SanitizedField)
}

type URLValidationResult struct {
	IsValid    bool
	ParsedURL  *url.URL
	Error      error
	IsInternal bool
}

type URLValidator interface {
	Validate(rawURL string) (*URLValidationResult, error)
	IsAllowedScheme(scheme string) bool
	IsInternalIP(ip string) bool
}

var AllowedHTTPSchemes = map[string]bool{
	"https": true,
	"http":  true,
}

var InternalIPPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^127\.`),                        // localhost
	regexp.MustCompile(`^10\.`),                         // private class A
	regexp.MustCompile(`^172\.(1[6-9]|2[0-9]|3[01])\.`), // private class B
	regexp.MustCompile(`^192\.168\.`),                   // private class C
	regexp.MustCompile(`^169\.254\.`),                   // link-local
	regexp.MustCompile(`^0\.`),                          // current network
	regexp.MustCompile(`^::1$`),                         // IPv6 localhost
	regexp.MustCompile(`^fc00:`),                        // IPv6 unique local
	regexp.MustCompile(`^fe80:`),                        // IPv6 link-local
}

type SecureHTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	Post(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Put(ctx context.Context, url string, body []byte, headers map[string]string) ([]byte, error)
	Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error)
	SetURLValidator(validator URLValidator)
}
