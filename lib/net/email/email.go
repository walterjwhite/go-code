package email

import (

	"bytes"
	"fmt"
	"github.com/emersion/go-message/mail"
	"time"
)


type EmailAccount struct {
	Username string
	Password string
	Domain   string

	EmailAddress *mail.Address

	ImapServer *EmailServer
	SmtpServer *EmailServer
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

func (e *EmailAccount) String() string {
	if e == nil {
		return "<nil EmailAccount>"
	}
	masked := ""
	if e.Password != "" {
		masked = "********"
	}
	var addr string
	if e.EmailAddress != nil {
		addr = e.EmailAddress.String()
	}
	return fmt.Sprintf("EmailAccount{Username:%s, Password:%s, Domain:%s, EmailAddress:%s, ImapServer:%v, SmtpServer:%v}",
		e.Username, masked, e.Domain, addr, e.ImapServer, e.SmtpServer)
}

func (m *EmailMessage) String() string {
	if m == nil {
		return "<nil EmailMessage>"
	}
	from := ""
	if m.From != nil {
		from = m.From.String()
	}
	to := make([]string, 0, len(m.To))
	for _, a := range m.To {
		if a != nil {
			to = append(to, a.String())
		}
	}
	cc := make([]string, 0, len(m.Cc))
	for _, a := range m.Cc {
		if a != nil {
			cc = append(cc, a.String())
		}
	}
	bcc := make([]string, 0, len(m.Bcc))
	for _, a := range m.Bcc {
		if a != nil {
			bcc = append(bcc, a.String())
		}
	}
	return fmt.Sprintf("EmailMessage{From:%s, To:%v, Cc:%v, Bcc:%v, Subject:%q, DateSent:%s, MessageId:%s}",
		from, to, cc, bcc, m.Subject, m.DateSent.Format(time.RFC3339), m.MessageId)
}


