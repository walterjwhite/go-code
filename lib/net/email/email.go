package email

import (
	"crypto/tls"
	"crypto/x509"

	"time"

	"bytes"
	"github.com/emersion/go-message/mail"
)

var (
	UserTlsConfig tls.Config
	rootCertPool  = x509.NewCertPool()
)

type EmailSenderAccount struct {
	Username string
	Password string
	Domain   string

	EmailAddress *mail.Address

	ImapServer         *EmailServer
	SmtpServer         *EmailServer
	Certificates       []string
	InsecureSkipVerify bool
}

type EmailServer struct {
	Host string
	Port int
}

type EmailAttachment struct {
	Name string
	Data *bytes.Buffer
}

type EmailMessage struct {
	From *mail.Address

	To  []*mail.Address
	Cc  []*mail.Address
	Bcc []*mail.Address

	Subject string
	Body    string

	DateSent time.Time

	Attachments []*EmailAttachment

	MessageId      string
	ConversationId string
}

func (e *EmailSenderAccount) Initialize() {
	e.addCerts()

	//UserTlsConfig = tls.Config{InsecureSkipVerify: e.InsecureSkipVerify, ServerName: e.Server.Host, RootCAs: rootCertPool}
}
