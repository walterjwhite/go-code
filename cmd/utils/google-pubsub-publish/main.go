package main

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/utils/remote/exec"

	"flag"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
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

	command = flag.String("command", "", "command to execute remotely")
)

func init() {
	application.ConfigureWithProperties(googleConf)
	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*command) == 0 {
		logging.Panic(errors.New("expecting command to be non-empty"))
	}

	r := exec.RemoteExec{Command: *command}
	r.Args = flag.Args()

	session := google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Publish(googleConf.TopicName, r)
}
