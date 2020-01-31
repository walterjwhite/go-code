package main

import (
	"flag"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/email"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
	"io"
	"os"
	"path/filepath"
)

var (
	//serverAddressFlag = flag.String("ServerAddress", "imap.gmail.com:993", "Server Address")
	//logoutTimeoutFlag = flag.Int("LogoutTimeout", 30, "Logout Timeout in seconds")
	usernameFlag          = flag.String("Username", "", "Username")
	passwordFlag          = flag.String("Password", "", "Password")
	sourceFolderFlag      = flag.String("SourceFolder", "INBOX", "Folder Name")
	destinationFolderFlag = flag.String("DestinationFolder", "Trash", "Folder Name")

	emailStorePath      = flag.String("EmailStorePath", "/tmp/email", "Email Store Path")
	attachmentStorePath = flag.String("AttachmentStorePath", "/tmp/attachments", "Email Attachment Store Path")
)

type moveConfiguration struct {
	Destination  string
	EmailSession *email.EmailSession
}

func init() {
	application.Configure()

	_initDirectory(*emailStorePath)
	_initDirectory(*attachmentStorePath)
}

func _initDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logging.Panic(os.MkdirAll(path, os.ModePerm))
	}
}

func main() {
	emailSenderAccount := &email.EmailSenderAccount{Username: *usernameFlag,
		Password:   *passwordFlag,
		ImapServer: &email.EmailServer{Host: "imap.gmail.com", Port: 993}}

	emailSession := emailSenderAccount.Connect()
	moveConfiguration := &moveConfiguration{Destination: *destinationFolderFlag, EmailSession: emailSession}

	// first, clear out the folder
	emailSession.Read(*sourceFolderFlag, moveConfiguration.processMessage, false)

	// then, read async to clear out the folder
	emailSession.ReadAsync(*sourceFolderFlag, moveConfiguration.processMessage, false)
}

func (c *moveConfiguration) processMessage(msg *imap.Message) {
	// convert to email message
	emailMessage := email.Process(msg)

	// TODO: make these actions configurable (store email, move)

	// store on disk
	filename := filepath.Join(*emailStorePath, fmt.Sprintf("%v", msg.Uid))
	yamlhelper.Write(emailMessage, filename)

	// store attachments on disk
	emailAttachmentDirectory := filepath.Join(*attachmentStorePath, fmt.Sprintf("%v", msg.Uid))

	for _, emailAttachment := range emailMessage.Attachments {
		writeAttachment(emailMessage, emailAttachmentDirectory, emailAttachment)
	}

	// move to folder
	c.EmailSession.Move(msg, c.Destination)

	// move message on disk to its final location (based on what rules it matches)
}

func writeAttachment(emailMessage *email.EmailMessage, emailAttachmentDirectory string, emailAttachment *email.EmailAttachment) {
	emailAttachmentFilename := filepath.Join(*attachmentStorePath, emailAttachment.Name)

	_initDirectory(emailAttachmentDirectory)

	// Create file with attachment name
	file, err := os.Create(emailAttachmentFilename)
	logging.Panic(err)

	// using io.Copy instead of io.ReadAll to avoid insufficient memory issues
	size, err := io.Copy(file, emailAttachment.Data)
	logging.Panic(err)

	log.Info().Msgf("Saved attachment: %v -> %v", size, file)
}
