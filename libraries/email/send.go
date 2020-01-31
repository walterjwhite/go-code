package email

import (
	"github.com/emersion/go-message/mail"
	sendmail "github.com/go-mail/gomail"
	"github.com/walterjwhite/go-application/libraries/logging"
)

func (e *EmailSenderAccount) Send(emailMessage *EmailMessage) {
	m := sendmail.NewMessage()

	setEmailAddressHeader(m, "From", e.EmailAddress)

	setEmailAddressHeader(m, "To", emailMessage.To...)
	setEmailAddressHeader(m, "Cc", emailMessage.Cc...)
	setEmailAddressHeader(m, "Bcc", emailMessage.Bcc...)

	m.SetHeader("Subject", emailMessage.Subject)
	m.SetBody("text/html", emailMessage.Body)

	d := sendmail.NewDialer(e.SmtpServer.Host, e.SmtpServer.Port, e.Username, e.Password)
	d.SSL = true
	//d.TLSConfig = &UserTlsConfig
	logging.Panic(d.DialAndSend(m))
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
