package daily_activity

import (
	"bytes"
	"time"

	"encoding/csv"
	"fmt"
	"net/mail"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/write"
)

func (c *Conf) sendEmail(body string) error {
	if c.emailAccount == nil {
		return fmt.Errorf("email account is nil")
	}

	msg := buildEmailMessage(c.emailAccount, body)
	log.Warn().Msgf("sending email to %v", c.emailAccount.EmailAddress)
	log.Warn().Msgf("email %v", msg)
	return write.Send(c.emailAccount, msg)
}

func buildEmailMessage(from *email.EmailAccount, body string) *email.EmailMessage {
	subject := fmt.Sprintf("daily activity - %s", time.Now().Format("2006/01/02"))
	return &email.EmailMessage{From: from.EmailAddress, To: []*mail.Address{from.EmailAddress}, Subject: subject, Body: body}
}

func generateCSVBody(cols []string, records []map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(cols); err != nil {
		return "", err
	}
	for _, r := range records {
		if err := w.Write(convertRecordToStrings(cols, r)); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
