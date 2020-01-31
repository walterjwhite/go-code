package main

import (
	"flag"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	usernameFlag = flag.String("Username", "", "Username")
	passwordFlag = flag.String("Password", "", "Password")
)

func init() {
	application.Configure()
}

func main() {
	log.Warn().Msg("This example does NOT work")

	log.Info().Msg("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	//c, err := client.Dial("imap.gmail.com:993")
	logging.Panic(err)
	log.Info().Msg("Connected")

	// Don't forget to logout
	defer logging.Panic(c.Logout())

	// Login
	err = c.Login(*usernameFlag, *passwordFlag)
	logging.Panic(err)
	log.Info().Msg("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Info().Msg("Mailboxes:")
	for m := range mailboxes {
		log.Info().Msg("* " + m.Name)
	}

	err = <-done
	logging.Panic(err)

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	logging.Panic(err)
	log.Info().Msgf("Flags for INBOX: %v", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Info().Msg("Last 4 messages:")
	for msg := range messages {
		log.Info().Msg("* " + msg.Envelope.Subject)
	}

	err = <-done
	logging.Panic(err)

	log.Info().Msg("Done!")
}
