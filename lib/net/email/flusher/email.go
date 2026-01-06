package flusher

import (
	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/write"
	"time"
)

type EmailFlusher struct {
	Account         *email.EmailAccount
	SubjectTemplate string
}

func (f *EmailFlusher) Flush(b []byte) error {
	if f == nil || f.Account == nil {
		return nil
	}

	err := write.Send(f.Account, f.toMessage(b))
	if err != nil {
		return err
	}

	log.Info().Int("bytes", len(b)).Msg("flushed buffer via email")
	return nil
}

func (f *EmailFlusher) toMessage(b []byte) *email.EmailMessage {
	return &email.EmailMessage{
		From:    f.Account.EmailAddress,
		To:      []*mail.Address{f.Account.EmailAddress},
		Subject: f.SubjectTemplate + " - " + time.Now().Format(time.RFC3339),
		Body:    string(b),
	}
}
