package email

import (
	"bytes"
	"errors"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io"
)

func Process(msg *imap.Message) *EmailMessage {

	emailMessage := &EmailMessage{}

	var section imap.BodySectionName
	r := msg.GetBody(&section)
	if r == nil {
		logging.Panic(errors.New("no message body returned"))
	}

	mr, err := mail.CreateReader(r)
	logging.Panic(err)

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
		logging.Panic(err)

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			b, _ := io.ReadAll(p.Body)
			emailMessage.Body = string(b)

			log.Info().Msgf("header: %v", p.Header)
			log.Info().Msgf("body: %v", emailMessage.Body[0:4])

			handleInlineAttachments(msg, p, h, emailMessage)

		case *mail.AttachmentHeader:
			handleAttachment(msg, p, h, emailMessage)
		}
	}

	return emailMessage
}

func handleInlineAttachments(msg *imap.Message, p *mail.Part, h *mail.InlineHeader, emailMessage *EmailMessage) {
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

func handleAttachment(msg *imap.Message, p *mail.Part, h *mail.AttachmentHeader, emailMessage *EmailMessage) {
	filename, _ := h.Filename()

	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, p.Body)
	logging.Panic(err)

	emailMessage.Attachments = append(emailMessage.Attachments, &EmailAttachment{Data: buffer, Name: filename})
}
