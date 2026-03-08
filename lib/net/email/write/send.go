package write

import (
	"crypto/tls"
	"fmt"
	stdmail "net/mail"
	"regexp"

	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
	gomail "gopkg.in/gomail.v2"

	"os"
	"strings"
)

var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var gomailDialAndSend = func(dialer *gomail.Dialer, m ...*gomail.Message) error {
	return dialer.DialAndSend(m...)
}

func validateEmailAddress(addr string) error {
	if addr == "" {
		return fmt.Errorf("email address cannot be empty")
	}

	if !emailPattern.MatchString(addr) {
		return fmt.Errorf("invalid email address format: %s", addr)
	}

	_, err := stdmail.ParseAddress(addr)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

func validateEmailAccount(e *email.EmailAccount) error {
	if e == nil {
		return fmt.Errorf("email account cannot be nil")
	}

	if e.SmtpServer == nil {
		return fmt.Errorf("SMTP server configuration cannot be nil")
	}

	if e.SmtpServer.Host == "" {
		return fmt.Errorf("SMTP server host cannot be empty")
	}

	if e.SmtpServer.Port <= 0 || e.SmtpServer.Port > 65535 {
		return fmt.Errorf("SMTP server port must be between 1 and 65535")
	}

	if e.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	return nil
}

func Send(e *email.EmailAccount, emailMessage *email.EmailMessage) error {
	if err := validateEmailAccount(e); err != nil {
		return fmt.Errorf("invalid email account configuration: %w", err)
	}

	if err := validateEmailMessage(emailMessage); err != nil {
		return fmt.Errorf("invalid email message: %w", err)
	}

	m, attachmentFilenames := buildMessage(e.EmailAddress, emailMessage)
	defer cleanupAttachments(attachmentFilenames)

	d := gomail.NewDialer(e.SmtpServer.Host, e.SmtpServer.Port, e.Username, e.Password)
	d.TLSConfig = &tls.Config{
		ServerName:         e.SmtpServer.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}
	return gomailDialAndSend(d, m)
}

func validateEmailMessage(m *email.EmailMessage) error {
	if m == nil {
		return fmt.Errorf("email message cannot be nil")
	}

	if m.From == nil {
		return fmt.Errorf("sender address (From) cannot be nil")
	}

	if err := validateEmailAddress(m.From.Address); err != nil {
		return fmt.Errorf("invalid sender address: %w", err)
	}

	if len(m.To) == 0 && len(m.Cc) == 0 && len(m.Bcc) == 0 {
		return fmt.Errorf("at least one recipient (To, Cc, or Bcc) must be specified")
	}

	for i, addr := range m.To {
		if addr == nil {
			return fmt.Errorf("recipient To[%d] cannot be nil", i)
		}
		if err := validateEmailAddress(addr.Address); err != nil {
			return fmt.Errorf("invalid recipient To[%d] address: %w", i, err)
		}
	}

	for i, addr := range m.Cc {
		if addr == nil {
			return fmt.Errorf("recipient Cc[%d] cannot be nil", i)
		}
		if err := validateEmailAddress(addr.Address); err != nil {
			return fmt.Errorf("invalid recipient Cc[%d] address: %w", i, err)
		}
	}

	for i, addr := range m.Bcc {
		if addr == nil {
			return fmt.Errorf("recipient Bcc[%d] cannot be nil", i)
		}
		if err := validateEmailAddress(addr.Address); err != nil {
			return fmt.Errorf("invalid recipient Bcc[%d] address: %w", i, err)
		}
	}

	return nil
}

func buildMessage(emailAddress *mail.Address, e *email.EmailMessage) (*gomail.Message, []string) {
	m := gomail.NewMessage()

	setHeader(m, "From", strings.TrimSpace(emailAddress.Address))

	setHeader(m, "To", addrsValToStrings(e.To)...)
	setHeader(m, "Cc", addrsValToStrings(e.Cc)...)
	setHeader(m, "Bcc", addrsValToStrings(e.Bcc)...)

	attachmentFilenames := addAttachments(e, m)

	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/plain", e.Body)

	return m, attachmentFilenames
}

var osCreateTemp = os.CreateTemp
var fileWrite = func(f *os.File, b []byte) (n int, err error) { return f.Write(b) }

const MaxAttachmentSize = 25 * 1024 * 1024

var allowedAttachmentExtensions = map[string]bool{
	".txt": true, ".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".zip": true, ".csv": true, ".json": true, ".xml": true,
}

func addAttachments(e *email.EmailMessage, m *gomail.Message) []string {
	attachmentFilenames := make([]string, 0)
	if len(e.Attachments) > 0 {
		for _, attachment := range e.Attachments {
			dataSize := len(attachment.Data.Bytes())
			if dataSize > MaxAttachmentSize {
				log.Error().Msgf("attachment %s exceeds maximum size of %d bytes", attachment.Name, MaxAttachmentSize)
				continue
			}

			if !isAllowedAttachmentType(attachment.Name) {
				log.Error().Msgf("attachment %s has disallowed file type", attachment.Name)
				continue
			}

			tmpFile, err := osCreateTemp("", "attach-*")
			if err != nil {
				log.Error().Err(err).Msg("failed to create temp file for attachment")
				continue
			}

			name := tmpFile.Name()
			if err := tmpFile.Chmod(0600); err != nil {
				log.Error().Err(err).Msg("failed to set restrictive permissions on temp file")
				_ = tmpFile.Close()
				_ = os.Remove(name)
				continue
			}

			log.Debug().Msgf("attachment size: %v", dataSize)
			if _, err := fileWrite(tmpFile, attachment.Data.Bytes()); err != nil {
				log.Error().Err(err).Msg("failed to write attachment to temp file")
				_ = tmpFile.Close()
				_ = os.Remove(name)
				continue
			}

			if err := tmpFile.Close(); err != nil {
				log.Error().Err(err).Msg("failed to close temp file for attachment")
				_ = os.Remove(name)
				continue
			}

			m.Attach(name)
			attachmentFilenames = append(attachmentFilenames, name)
		}
	}

	return attachmentFilenames
}

func isAllowedAttachmentType(filename string) bool {
	if filename == "" {
		return false
	}

	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	return allowedAttachmentExtensions[ext]
}

func cleanupAttachments(attachmentFilenames []string) {
	if len(attachmentFilenames) == 0 {
		return
	}

	for _, attachmentFilename := range attachmentFilenames {
		if err := os.Remove(attachmentFilename); err != nil {
			log.Warn().Msg("failed to cleanup temporary attachment file")
		}
	}
}


func addrsValToStrings(addrs []*mail.Address) []string {
	out := make([]string, 0, len(addrs))
	for i := range addrs {
		cleanAddr := sanitizeHeaderValue(strings.TrimSpace(addrs[i].Address))
		if cleanAddr != "" {
			out = append(out, cleanAddr)
		}
	}

	return out
}

func sanitizeHeaderValue(value string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		}
		return r
	}, value)
}

func setHeader(m *gomail.Message, headerName string, value ...string) {
	sanitizedValue := make([]string, 0, len(value))
	for _, v := range value {
		clean := sanitizeHeaderValue(v)
		if clean != "" {
			sanitizedValue = append(sanitizedValue, clean)
		}
	}

	if len(sanitizedValue) > 0 && strings.Join(sanitizedValue, "") != "" {
		log.Debug().Msgf("setting header: %s -> %v", headerName, strings.Join(sanitizedValue, ", "))
		m.SetHeader(headerName, sanitizedValue...)
	}
}
