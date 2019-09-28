package email

import (
	gomail "github.com/go-mail/gomail"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (e *EmailSenderAccount) Send(emailMessage EmailMessage) {
	m := gomail.NewMessage()

	m.SetHeader("From", e.EmailAddress)

	setHeader(m, "To", emailMessage.To...)
	setHeader(m, "Cc", emailMessage.Cc...)
	setHeader(m, "Bcc", emailMessage.Bcc...)

	m.SetHeader("Subject", emailMessage.Subject)
	m.SetBody("text/html", emailMessage.Body)

	d := gomail.NewDialer(e.Server.Host, e.Server.Port, e.Username, e.Password)
	d.SSL = false
	d.TLSConfig = &UserTlsConfig
	logging.Panic(d.DialAndSend(m))
}

func setHeader(m *gomail.Message, headerName string, value ...string) {
	if len(value) > 0 {
		m.SetHeader(headerName, value...)
	}
}
