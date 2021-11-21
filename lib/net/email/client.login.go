package email

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type EmailSession struct {
	EmailSenderAccount *EmailSenderAccount

	client *client.Client
}

func (a *EmailSenderAccount) Connect() *EmailSession {
	emailSession := &EmailSession{EmailSenderAccount: a}

	c, err := client.DialTLS(fmt.Sprintf("%v:%v", a.ImapServer.Host, a.ImapServer.Port), nil)
	logging.Panic(err)

	logging.Panic(c.Login(a.Username, a.Password))

	emailSession.client = c
	return emailSession
}

func (s *EmailSession) Clone() *EmailSession {
	return s.EmailSenderAccount.Connect()
}

func (s *EmailSession) Close() {
	logging.Panic(s.client.Logout())
}
