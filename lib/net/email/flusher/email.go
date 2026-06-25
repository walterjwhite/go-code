package flusher

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/write"
)

type EmailFlusher struct {
	Account         *email.EmailAccount
	SubjectTemplate string
}

func (f *EmailFlusher) Flush(b []byte) error {
	if f == nil || f.Account == nil {
		return errors.New("flusher or account is nil")
	}

	if len(b) == 0 {
		return errors.New("cannot flush empty buffer")
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
		Subject: f.SubjectTemplate + " - " + time.Now().Format(time.RFC3339) + " - " + generateSecureID(),
		Body:    string(b),
	}
}

func generateSecureID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405.000000")
	}
	return hex.EncodeToString(bytes)
}
