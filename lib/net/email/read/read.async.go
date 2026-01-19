package write

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/emersion/go-imap-idle"
)

func (s *EmailSession) ReadAsync(folderName string, function func(msg *imap.Message), incrementIndex bool) error {
	_, err := s.client.Select(folderName, false)
	if err != nil {
		return err
	}

	idleClient := idle.NewClient(s.client)

	updates := make(chan client.Update)
	s.client.Updates = updates

	done := make(chan error, 1)
	go func() {
		done <- idleClient.IdleWithFallback(nil, 0)
	}()

	for {
		select {
		case update := <-updates:
			log.Info().Msgf("New update: %v", update)

			mailboxUpdate, ok := update.(*client.MailboxUpdate)
			if ok {
				err = s.readFolder(mailboxUpdate, function, incrementIndex)
				if err != nil {
					return err
				}
			}
		case err := <-done:
			logging.Error(err)
			log.Info().Msg("Not idling anymore")
			return nil
		}
	}

}

func (s *EmailSession) readFolder(mailboxUpdate *client.MailboxUpdate, function func(msg *imap.Message), incrementIndex bool) error {
	log.Info().Msgf("Mailbox Update: %v %v", mailboxUpdate.Mailbox.Name, mailboxUpdate.Mailbox.Items)

	sClone, err := s.clone()
	if err != nil {
		return err
	}

	defer sClone.Close()

	return sClone.Read(mailboxUpdate.Mailbox.Name, function, incrementIndex)
}

func (s *EmailSession) clone() (*EmailSession, error) {
	ns, err := New(s.EmailAccount)
	if err != nil {
		return nil, err
	}

	return ns, nil
}
