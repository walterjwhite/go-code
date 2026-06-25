package write

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/net/email"
	"io"
)

func ImapMessageToEmailMessage(msg *imap.Message) (*email.EmailMessage, error) {
	emailMessage := &email.EmailMessage{}

	var section imap.BodySectionName
	r := msg.GetBody(&section)
	if r == nil {
		return nil, errors.New("no message body returned")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail reader: %w", err)
	}

	emailMessage.DateSent = msg.Envelope.Date
	emailMessage.Subject = msg.Envelope.Subject

	emailMessage.MessageId = msg.Envelope.MessageId
	emailMessage.ConversationId = msg.Envelope.InReplyTo

	if from, err := mr.Header.AddressList("From"); err == nil {
		emailMessage.From = from[0]
	}
	if to, err := mr.Header.AddressList("To"); err == nil {
		emailMessage.To = to
	}
	if bcc, err := mr.Header.AddressList("Bcc"); err == nil {
		emailMessage.Bcc = bcc
	}
	if cc, err := mr.Header.AddressList("Cc"); err == nil {
		emailMessage.Cc = cc
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read message part: %w", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, err := io.ReadAll(p.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read message body: %w", err)
			}
			emailMessage.Body = string(b)

			handleInlineAttachments(msg, p, h, emailMessage)

		case *mail.AttachmentHeader:
			err := handleAttachment(p, h, emailMessage)
			if err != nil {
				return nil, err
			}
		}
	}

	return emailMessage, nil
}

func handleInlineAttachments(msg *imap.Message, p *mail.Part, h *mail.InlineHeader, emailMessage *email.EmailMessage) {
	log.Info().Msgf("Inline header: %v", h)
	/*
		headerDisplay, cdparams, err := h.ContentDisposition()
		if err != nil {
			log.Warn().Msgf("Error with inline header: %v", err)
		}

		log.Info().Msgf("headerDisplay: %v", headerDisplay, cdparams)

		headerType, params, err := h.ContentType()
		if err != nil {
			log.Warn().Msgf("Error with inline header: %v", err)
		}

		log.Info().Msgf("headerType: %v %v", headerType, params)
	*/
}

func handleAttachment(p *mail.Part, h *mail.AttachmentHeader, emailMessage *email.EmailMessage) error {
	filename, err := h.Filename()
	if err != nil {
		return fmt.Errorf("failed to get attachment filename: %w", err)
	}

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, p.Body)
	if err != nil {
		return fmt.Errorf("failed to read attachment data: %w", err)
	}

	emailMessage.Attachments = append(emailMessage.Attachments, &email.EmailAttachment{Data: buffer, Name: filename})
	return nil
}
