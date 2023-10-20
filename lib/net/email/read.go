package email

import (
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (s *EmailSession) Read(folderName string, function func(msg *imap.Message), incrementIndex bool) {
	mbox, err := s.client.Select(folderName, false)
	logging.Panic(err)

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

		logging.Panic(<-done)

		function(<-messageChannel)
		seqSet.Clear()

		if incrementIndex {
			i++
		}
	}
}

/*
	yaml.Write(emailMessage, "/tmp/email-message")

	// push to ES
	ces := elasticsearch.NewDefaultClient()
	b := ces.NewBatch(11, 5242880, 5*time.Second, 2)

	b.Index("main.email.1", emailMessage)

	b.Flush()

	log.Info().Msg("indexed email message")
*/
