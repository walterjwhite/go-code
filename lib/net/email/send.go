package email

import (
	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"

	sendmail "github.com/go-mail/gomail"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"os"
)

// NOTE: for gmail users, insecure access must be enabled for this to work
func (e *EmailSenderAccount) Send(emailMessage *EmailMessage) {
	m := sendmail.NewMessage()

	setEmailAddressHeader(m, "From", e.EmailAddress)

	setEmailAddressHeader(m, "To", emailMessage.To...)
	setEmailAddressHeader(m, "Cc", emailMessage.Cc...)
	setEmailAddressHeader(m, "Bcc", emailMessage.Bcc...)

	attachmentFilenames := addAttachments(emailMessage, m)

	m.SetHeader("Subject", emailMessage.Subject)
	// m.SetBody("text/html", emailMessage.Body)
	m.SetBody("text/plain", emailMessage.Body)

	log.Debug().Msgf("subject: %s", emailMessage.Subject)
	log.Debug().Msgf("body: %s", emailMessage.Body)

	d := sendmail.NewDialer(e.SmtpServer.Host, e.SmtpServer.Port, e.Username, e.Password)
	d.SSL = true
	//d.TLSConfig = &UserTlsConfig
	logging.Panic(d.DialAndSend(m))

	cleanupAttachments(attachmentFilenames)
}

func addAttachments(emailMessage *EmailMessage, m *sendmail.Message) []string {
	attachmentFilenames := make([]string, 0)
	if len(emailMessage.Attachments) > 0 {
		for _, attachment := range emailMessage.Attachments {
			tmpFile, err := os.CreateTemp(os.TempDir(), "*"+attachment.Name)
			logging.Panic(err)

			log.Debug().Msgf("attachment size: %v", len(attachment.Data.Bytes()))
			logging.Panic(os.WriteFile(tmpFile.Name(), attachment.Data.Bytes(), 0644))

			m.Attach(tmpFile.Name())
			attachmentFilenames = append(attachmentFilenames, tmpFile.Name())
		}
	}

	return attachmentFilenames
}

func cleanupAttachments(attachmentFilenames []string) {
	if len(attachmentFilenames) > 0 {
		for _, attachmentFilename := range attachmentFilenames {
			os.Remove(attachmentFilename)
		}
	}

}

func setEmailAddressHeader(m *sendmail.Message, headerName string, emailAddresses ...*mail.Address) {
	for _, emailAddress := range emailAddresses {
		setHeader(m, headerName, emailAddress.Address)
	}
}

func setHeader(m *sendmail.Message, headerName string, value ...string) {
	if len(value) > 0 {
		m.SetHeader(headerName, value...)
	}
}
