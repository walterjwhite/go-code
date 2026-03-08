package write

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	"regexp"
)

var folderNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._\-/\s]+$`)

func validateFolderName(folderName string) error {
	if folderName == "" {
		return fmt.Errorf("folder name cannot be empty")
	}
	if len(folderName) > 255 {
		return fmt.Errorf("folder name exceeds maximum length of 255 characters")
	}
	if !folderNamePattern.MatchString(folderName) {
		return fmt.Errorf("folder name contains invalid characters")
	}
	return nil
}

func (s *EmailSession) Read(folderName string, function func(msg *imap.Message), incrementIndex bool) error {
	if err := validateFolderName(folderName); err != nil {
		return fmt.Errorf("invalid folder name: %w", err)
	}

	mbox, err := s.client.Select(folderName, false)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v Messages.", mbox.Messages)

	seqSet := new(imap.SeqSet)

	var section imap.BodySectionName
	var i uint32
	for i = 1; i <= mbox.Messages; {
		messageChannel := make(chan *imap.Message, 1)
		done := make(chan error, 1)

		items := []imap.FetchItem{section.FetchItem(), imap.FetchUid, imap.FetchEnvelope}

		log.Info().Msgf("i:%v", i)
		seqSet.AddNum(i)

		go func() {
			done <- s.client.Fetch(seqSet, items, messageChannel)
		}()

		err = <-done
		if err != nil {
			close(messageChannel)
			return fmt.Errorf("fetch failed: %w", err)
		}

		msg := <-messageChannel
		close(messageChannel)
		function(msg)
		seqSet.Clear()

		if incrementIndex {
			i++
		}
	}

	return nil
}

/*
	yaml.Write(emailMessage, "/tmp/email-message")

	ces := elasticsearch.NewDefaultClient()
	b := ces.NewBatch(11, 5242880, 5*time.Second, 2)

	b.Index("main.email.1", emailMessage)

	b.Flush()

	log.Info().Msg("indexed email message")
*/
