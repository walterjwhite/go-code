package write

import (
	"fmt"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/email"
)

type EmailSession struct {
	EmailAccount *email.EmailAccount

	client *client.Client
}

func New(a *email.EmailAccount) (*EmailSession, error) {
	c, err := client.DialTLS(fmt.Sprintf("%v:%v", a.ImapServer.Host, a.ImapServer.Port), nil)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("credentials for user: %v (password redacted)", a.Username)
	err = c.Login(a.Username, a.Password)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("successfully logged in")

	return &EmailSession{EmailAccount: a, client: c}, nil
}

func (s *EmailSession) Close() {
	if s.client == nil {
		return
	}

	logging.Warn(s.client.Logout(), "EmailSession.Close")
	s.client = nil
}
