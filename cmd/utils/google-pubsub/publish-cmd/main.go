package main

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/exec"
	"github.com/walterjwhite/go-code/lib/net/google"

	"flag"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"strings"
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

	functionName = flag.String("functionName", "", "function to execute remotely")
	arguments    = flag.String("arguments", "", "arguments to pass functionName, optional")
)

func init() {
	application.Configure(googleConf)
	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*functionName) == 0 {
		logging.Panic(errors.New("expecting command to be non-empty"))
	}

	c := exec.Cmd{FunctionName: *functionName}
	if len(*arguments) != 0 {
		c.Args = strings.Fields(*arguments)
	}

	session := google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Publish(googleConf.TopicName, c)
}
