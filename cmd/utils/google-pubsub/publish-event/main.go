package main

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"

	"flag"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"os"
)

type PublisherConfiguration struct {
	CredentialsFile string
	ProjectId       string

	TopicName string

	EncryptionKeyFilename string
}

var (
	googleConf = &PublisherConfiguration{}
	aesConf    = &aes.Configuration{}

	message     = flag.String("message", "", "message to publish, takes precedence of file")
	messageFile = flag.String("file", "", "message file to use for contents of message")
)

func init() {
	application.Configure(googleConf)
	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*message) == 0 {
		if len(*messageFile) == 0 {
			logging.Panic(errors.New("expecting command to be non-empty"))
		}

		fileContents, err := os.ReadFile(*messageFile)
		logging.Panic(err)

		messageContents := string(fileContents)
		message = &messageContents
	}

	session := google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Publish(googleConf.TopicName, message)
}
