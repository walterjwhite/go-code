package main

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"

	"flag"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"strings"
	"os/exec"
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

	command   = flag.String("command", "", "command to execute remotely")
	arguments = flag.String("arguments", "", "arguments to pass cmd, optional")
)

func init() {
	application.Configure(googleConf)
	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*command) == 0 {
		logging.Panic(errors.New("expecting command to be non-empty"))
	}

	r := exec.Command(*command)
	if len(*arguments) != 0 {
		r.Args = strings.Fields(*arguments)
	}

	session := google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Publish(googleConf.TopicName, r)
}
