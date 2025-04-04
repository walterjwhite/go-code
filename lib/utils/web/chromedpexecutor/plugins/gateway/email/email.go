package email

import (
	"context"
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	emaill "github.com/walterjwhite/go-code/lib/net/email"
	"regexp"
)

type Provider struct {
	EmailSenderAccount *emaill.EmailSenderAccount

	emailSession *emaill.EmailSession

	channel chan *string
}

const (
	TOKEN_REGEXP = "^token: [0-9]{6}$"
	INBOX_FOLDER = "INBOX"
)

func (p *Provider) Get(ctx context.Context) string {
	p.channel = make(chan *string, 1)

	log.Info().Msgf("Connecting to: %v @ %v:%v", p.EmailSenderAccount.EmailAddress, p.EmailSenderAccount.ImapServer.Host, p.EmailSenderAccount.ImapServer.Port)

	p.emailSession = p.EmailSenderAccount.Connect()
	go p.emailSession.ReadAsync(INBOX_FOLDER, p.onNewMessage, true)

	return *<-p.channel
}

func (p *Provider) onNewMessage(msg *imap.Message) {
	if isMessageToken(msg) {
		log.Info().Msgf("is token: %v", msg.Envelope.Subject)
		p.onTokenReceived(msg)
	} else {
		log.Debug().Msgf("NOT token: %v", msg.Envelope.Subject)
	}
}

func isMessageToken(msg *imap.Message) bool {
	r := regexp.MustCompile(TOKEN_REGEXP)

	return r.MatchString(msg.Envelope.Subject)
}

func (p *Provider) onTokenReceived(msg *imap.Message) {
	t := getToken(msg.Envelope.Subject)
	log.Debug().Msgf("token being sent to channel: %v", t)

	p.channel <- &t

	log.Debug().Msgf("token delivered to channel: %v", t)
}

func getToken(subject string) string {
	return subject[7:]
}
