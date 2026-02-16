package write

import (
	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/email"
	gomail "gopkg.in/gomail.v2"

	"os"
	"strings"
)

var gomailDialAndSend = func(dialer *gomail.Dialer, m ...*gomail.Message) error {
	return dialer.DialAndSend(m...)
}

func Send(e *email.EmailAccount, emailMessage *email.EmailMessage) error {
	m, attachmentFilenames := buildMessage(e.EmailAddress, emailMessage)
	defer cleanupAttachments(attachmentFilenames)

	d := gomail.NewDialer(e.SmtpServer.Host, e.SmtpServer.Port, e.Username, e.Password)
	return gomailDialAndSend(d, m)
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

	log.Warn().Msgf("subject: %s", e.Subject)
	log.Warn().Msgf("body: %s", e.Body)

	return m, attachmentFilenames
}

var osCreateTemp = os.CreateTemp
var fileWrite = func(f *os.File, b []byte) (n int, err error) { return f.Write(b) }

func addAttachments(e *email.EmailMessage, m *gomail.Message) []string {
	attachmentFilenames := make([]string, 0)
	if len(e.Attachments) > 0 {
		for _, attachment := range e.Attachments {
			tmpFile, err := osCreateTemp("", "attach-*")
			if err != nil {
				log.Error().Err(err).Msg("failed to create temp file for attachment")
				continue
			}

			name := tmpFile.Name()
			defer func() {
				_ = tmpFile.Close()
			}()

			log.Debug().Msgf("attachment size: %v", len(attachment.Data.Bytes()))
			if _, err := fileWrite(tmpFile, attachment.Data.Bytes()); err != nil {
				log.Error().Err(err).Msg("failed to write attachment to temp file")
				_ = tmpFile.Close()
				_ = os.Remove(name)
				continue
			}

			_ = tmpFile.Chmod(0600)
			_ = tmpFile.Close()

			m.Attach(name)
			attachmentFilenames = append(attachmentFilenames, name)
		}
	}

	return attachmentFilenames
}

func cleanupAttachments(attachmentFilenames []string) {
	if len(attachmentFilenames) == 0 {
		return
	}

	for _, attachmentFilename := range attachmentFilenames {
		logging.Warn(os.Remove(attachmentFilename), "cleanupAttachments")
	}
}


func addrsValToStrings(addrs []*mail.Address) []string {
	out := make([]string, 0, len(addrs))
	for i := range addrs {

		out = append(out, strings.TrimSpace(addrs[i].Address))
	}

	return out
}

func setHeader(m *gomail.Message, headerName string, value ...string) {
	if len(value) > 0 && strings.Join(value, "") != "" {
		log.Debug().Msgf("setting header: %s -> %v", headerName, strings.Join(value, ", "))
		m.SetHeader(headerName, value...)
	}
}
