package main

import (
	"flag"
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/email"
)

var (
	//serverAddressFlag = flag.String("ServerAddress", "imap.gmail.com:993", "Server Address")
	//logoutTimeoutFlag = flag.Int("LogoutTimeout", 30, "Logout Timeout in seconds")
	usernameFlag = flag.String("Username", "", "Username")
	passwordFlag = flag.String("Password", "", "Password")
	folderFlag   = flag.String("Folder", "INBOX", "Folder Name")
)

type moveConfiguration struct {
	Destination  string
	EmailSession *email.EmailSession
}

func init() {
	application.Configure()
}

func main() {
	emailSenderAccount := &email.EmailSenderAccount{Username: *usernameFlag,
		Password:   *passwordFlag,
		ImapServer: &email.EmailServer{Host: "imap.gmail.com", Port: 993}}

	emailSession := emailSenderAccount.Connect()
	moveConfiguration := &moveConfiguration{Destination: "Trash", EmailSession: emailSession}

	emailSession.Read(*folderFlag, moveConfiguration.moveMessage, false)
}

func logMessage(msg *imap.Message) {
	log.Info().Msgf("read: %v", *email.Process(msg))
}

func (c *moveConfiguration) moveMessage(msg *imap.Message) {
	c.EmailSession.Move(msg, c.Destination)
}
