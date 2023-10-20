package email

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/emersion/go-imap-idle"
)

func (s *EmailSession) ReadAsync(folderName string, function func(msg *imap.Message), incrementIndex bool) {
	_, err := s.client.Select(folderName, false)
	logging.Panic(err)

	idleClient := idle.NewClient(s.client)

	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	s.client.Updates = updates

	// Start idling
	done := make(chan error, 1)
	go func() {
		done <- idleClient.IdleWithFallback(nil, 0)
	}()

	// Listen for updates
	for {
		select {
		case update := <-updates:
			log.Info().Msgf("New update: %v", update)

			mailboxUpdate, ok := update.(*client.MailboxUpdate)
			if ok {
				s.readFolder(mailboxUpdate, function, incrementIndex)
			}
		case err := <-done:
			logging.Panic(err)
			log.Info().Msg("Not idling anymore")
			return
		}
	}
}

func (s *EmailSession) readFolder(mailboxUpdate *client.MailboxUpdate, function func(msg *imap.Message), incrementIndex bool) {
	log.Info().Msgf("Mailbox Update: %v %v", mailboxUpdate.Mailbox.Name, mailboxUpdate.Mailbox.Items)

	sClone := s.Clone()
	defer sClone.Close()

	sClone.Read(mailboxUpdate.Mailbox.Name, function, incrementIndex)
}
