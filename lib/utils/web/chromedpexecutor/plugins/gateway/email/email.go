package email

import (
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	emaill "github.com/walterjwhite/go/lib/net/email"
	"regexp"
)

type Provider struct {
	EmailSenderAccount *emaill.EmailSenderAccount

	emailSession *emaill.EmailSession

	channel chan *string
}

// block until we receive the email
func (p *Provider) Get() string {
	p.channel = make(chan *string, 1)

	p.emailSession = p.EmailSenderAccount.Connect()
	// TODO: do not hard-code this folder
	go p.emailSession.ReadAsync("INBOX", p.onNewMessage, true)

	return *<-p.channel
}

func (p *Provider) onNewMessage(msg *imap.Message) {
	if isMessageToken(msg) {
		log.Debug().Msgf("is token: %v", msg.Envelope.Subject)
		p.onTokenReceived(msg)
	} else {
		log.Debug().Msgf("NOT token: %v", msg.Envelope.Subject)
	}
}

func isMessageToken(msg *imap.Message) bool {
	r := regexp.MustCompile("^token: [0-9]{6}$")

	return r.MatchString(msg.Envelope.Subject)
}

func (p *Provider) onTokenReceived(msg *imap.Message) {
	t := getToken(msg.Envelope.Subject)
	log.Debug().Msgf("token being sent to channel: %v", t)

	p.channel <- &t

	log.Debug().Msgf("token delivered to channel: %v", t)
	//defer p.emailSession.Delete(msg)
}

func getToken(subject string) string {
	return subject[7:]
}
