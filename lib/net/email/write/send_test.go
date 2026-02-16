package write

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/emersion/go-message/mail"
	"github.com/walterjwhite/go-code/lib/net/email"
	gomail "gopkg.in/gomail.v2"
)

func TestAddrsValToStrings(t *testing.T) {
	tests := []struct {
		name  string
		addrs []*mail.Address
		want  []string
	}{
		{
			name:  "single address",
			addrs: []*mail.Address{{Address: "test@example.com"}},
			want:  []string{"test@example.com"},
		},
		{
			name:  "multiple addresses",
			addrs: []*mail.Address{{Address: "test1@example.com"}, {Address: "test2@example.com"}},
			want:  []string{"test1@example.com", "test2@example.com"},
		},
		{
			name:  "address with leading/trailing spaces",
			addrs: []*mail.Address{{Address: "  test@example.com  "}},
			want:  []string{"test@example.com"},
		},
		{
			name:  "empty list",
			addrs: []*mail.Address{},
			want:  []string{},
		},
		{
			name:  "nil list",
			addrs: nil,
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addrsValToStrings(tt.addrs)
			if !compareStringSlices(got, tt.want) {
				t.Errorf("addrsValToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestSetHeader(t *testing.T) {
	m := gomail.NewMessage()

	setHeader(m, "Subject", "Test Subject")
	if m.GetHeader("Subject")[0] != "Test Subject" {
		t.Errorf("Expected subject 'Test Subject', got '%s'", m.GetHeader("Subject")[0])
	}

	setHeader(m, "To", "test1@example.com", "test2@example.com")
	if !compareStringSlices(m.GetHeader("To"), []string{"test1@example.com", "test2@example.com"}) {
		t.Errorf("Expected To headers ['test1@example.com', 'test2@example.com'], got %v", m.GetHeader("To"))
	}

	m = gomail.NewMessage() // Reset message
	setHeader(m, "X-Empty", "")
	if len(m.GetHeader("X-Empty")) != 0 {
		t.Errorf("Expected empty X-Empty header, got %v", m.GetHeader("X-Empty"))
	}
}

func TestAddAttachments(t *testing.T) {
	attachmentContent := "This is a test attachment content."
	attachmentBuffer := bytes.NewBufferString(attachmentContent)

	emailMessage := &email.EmailMessage{
		Attachments: []*email.EmailAttachment{
			{
				Data: attachmentBuffer,
				Name: "test_attachment.txt",
			},
		},
	}

	m := gomail.NewMessage()
	attachmentFilenames := addAttachments(emailMessage, m)

	if len(attachmentFilenames) != 1 {
		t.Fatalf("Expected 1 attachment filename, got %d", len(attachmentFilenames))
	}


	tempFilePath := attachmentFilenames[0]
	content, err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to read temporary attachment file: %v", err)
	}
	if string(content) != attachmentContent {
		t.Errorf("Temporary attachment file content mismatch. Got: '%s', Want: '%s'", string(content), attachmentContent)
	}

	defer cleanupAttachments(attachmentFilenames)
}

func TestAddAttachments_NoAttachments(t *testing.T) {
	emailMessage := &email.EmailMessage{
		Attachments: []*email.EmailAttachment{},
	}

	m := gomail.NewMessage()
	attachmentFilenames := addAttachments(emailMessage, m)

	if len(attachmentFilenames) != 0 {
		t.Errorf("Expected 0 attachment filenames, got %d", len(attachmentFilenames))
	}
}





func TestAddAttachments_CreateTempError(t *testing.T) {
	originalOsCreateTemp := osCreateTemp
	defer func() { osCreateTemp = originalOsCreateTemp }()

	osCreateTemp = func(dir, pattern string) (*os.File, error) {
		return nil, errors.New("mock create temp error")
	}

	emailMessage := &email.EmailMessage{
		Attachments: []*email.EmailAttachment{
			{
				Data: bytes.NewBufferString("test content"),
				Name: "test_attachment.txt",
			},
		},
	}

	m := gomail.NewMessage()
	attachmentFilenames := addAttachments(emailMessage, m)

	if len(attachmentFilenames) != 0 {
		t.Errorf("Expected 0 attachment filenames, got %d", len(attachmentFilenames))
	}
}

func TestAddAttachments_WriteError(t *testing.T) {
	originalOsCreateTemp := osCreateTemp
	originalFileWrite := fileWrite
	defer func() {
		osCreateTemp = originalOsCreateTemp
		fileWrite = originalFileWrite
	}()

	osCreateTemp = func(dir, pattern string) (*os.File, error) {
		return os.NewFile(0, "mock-file"), nil // Return a dummy os.File
	}

	fileWrite = func(f *os.File, b []byte) (n int, err error) {
		return 0, errors.New("mock write error")
	}

	emailMessage := &email.EmailMessage{
		Attachments: []*email.EmailAttachment{
			{
				Data: bytes.NewBufferString("test content"),
				Name: "test_attachment.txt",
			},
		},
	}

	m := gomail.NewMessage()
	attachmentFilenames := addAttachments(emailMessage, m)

	if len(attachmentFilenames) != 0 {
		t.Errorf("Expected 0 attachment filenames, got %d", len(attachmentFilenames))
	}
}

func TestCleanupAttachments(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-cleanup-*")
	if err != nil {
		t.Fatalf("Failed to create temp file for cleanup test: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file for cleanup test: %v", err)
	}

	tempFilePath := tmpFile.Name()

	_, err = os.Stat(tempFilePath)
	if os.IsNotExist(err) {
		t.Fatalf("Temporary file does not exist before cleanup: %s", tempFilePath)
	}

	cleanupAttachments([]string{tempFilePath})

	_, err = os.Stat(tempFilePath)
	if !os.IsNotExist(err) {
		t.Errorf("Temporary file was not cleaned up: %s", tempFilePath)
	}
}

func TestBuildMessage(t *testing.T) {
	senderAddress := &mail.Address{Address: "sender@example.com"}
	toAddress := []*mail.Address{{Address: "recipient@example.com"}}
	ccAddress := []*mail.Address{{Address: "cc@example.com"}}
	bccAddress := []*mail.Address{{Address: "bcc@example.com"}}

	emailMessage := &email.EmailMessage{
		To:      toAddress,
		Cc:      ccAddress,
		Bcc:     bccAddress,
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	m, attachmentFilenames := buildMessage(senderAddress, emailMessage)
	defer cleanupAttachments(attachmentFilenames) // Clean up any potential attachments

	if m.GetHeader("From")[0] != senderAddress.Address {
		t.Errorf("Expected From header '%s', got '%s'", senderAddress.Address, m.GetHeader("From")[0])
	}
	if !compareStringSlices(m.GetHeader("To"), addrsValToStrings(toAddress)) {
		t.Errorf("Expected To header '%v', got '%v'", addrsValToStrings(toAddress), m.GetHeader("To"))
	}
	if !compareStringSlices(m.GetHeader("Cc"), addrsValToStrings(ccAddress)) {
		t.Errorf("Expected Cc header '%v', got '%v'", addrsValToStrings(ccAddress), m.GetHeader("Cc"))
	}
	if m.GetHeader("Subject")[0] != emailMessage.Subject {
		t.Errorf("Expected Subject header '%s', got '%s'", emailMessage.Subject, m.GetHeader("Subject")[0])
	}

}

func TestSend(t *testing.T) {
	originalGomailDialAndSend := gomailDialAndSend
	defer func() { gomailDialAndSend = originalGomailDialAndSend }()

	t.Run("successful send", func(t *testing.T) {
		gomailDialAndSend = func(dialer *gomail.Dialer, m ...*gomail.Message) error {
			return nil // Simulate successful send
		}

		emailAccount := &email.EmailAccount{
			EmailAddress: &mail.Address{Address: "sender@example.com"},
			SmtpServer:   &email.EmailServer{Host: "smtp.example.com", Port: 587},
			Username:     "user",
			Password:     "pass",
		}
		emailMessage := &email.EmailMessage{
			To:      []*mail.Address{{Address: "recipient@example.com"}},
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := Send(emailAccount, emailMessage)
		if err != nil {
			t.Errorf("Send() failed unexpectedly: %v", err)
		}
	})

	t.Run("send with error", func(t *testing.T) {
		gomailDialAndSend = func(dialer *gomail.Dialer, m ...*gomail.Message) error {
			return errors.New("mock send error") // Simulate send error
		}

		emailAccount := &email.EmailAccount{
			EmailAddress: &mail.Address{Address: "sender@example.com"},
			SmtpServer:   &email.EmailServer{Host: "smtp.example.com", Port: 587},
			Username:     "user",
			Password:     "pass",
		}
		emailMessage := &email.EmailMessage{
			To:      []*mail.Address{{Address: "recipient@example.com"}},
			Subject: "Test Subject",
			Body:    "Test Body",
		}

		err := Send(emailAccount, emailMessage)
		if err == nil {
			t.Errorf("Send() expected an error, but got none")
		}
	})
}
