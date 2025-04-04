package main

import (
	"errors"

	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	"os/exec"
)

type SubscriberConfiguration struct {
	CredentialsFile string
	ProjectId       string

	TopicName        string
	SubscriptionName string

	ResponseTopicName     string
	EncryptionKeyFilename string
}

type Command struct {
	command    *string
	remoteExec exec.Command
}

var (
	googleConf = &SubscriberConfiguration{}
	aesConf    = &aes.Configuration{}
	session    *google.Session

	cmd = flag.String("cmd", "", "cmd to execute on receipt of pubsub message")
)

func init() {
	application.Configure(googleConf)

	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*cmd) == 0 {
		logging.Panic(errors.New("-cmd=<COMMAND>"))
	}

	c := &Command{command: cmd}

	session = google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Subscribe(googleConf.TopicName, googleConf.SubscriptionName, c)
	application.Wait()
}

func (c *Command) New() any {
	c.remoteExec = exec.Command("")
	return &c.remoteExec
}

func (c *Command) MessageDeserialized() {
	log.Info().Msgf("running: %s -> %s", c.remoteExec.Command, c.remoteExec.Args)
	c.remoteExec.Args = append([]string{c.remoteExec.Command}, c.remoteExec.Args...)

	ecmd := exec.Command(*c.command, c.remoteExec.Args...)
	output, err := ecmd.Output()

	status := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
			log.Warn().Msgf("Error running: %s (%s) -> %v", c.remoteExec.Command, c.remoteExec.Args, status)
		}
	}

	respond(status, string(output))
}

func (c *Command) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}

func respond(status int, output string) {
	log.Info().Msgf("status: %v, output: %s", status, output)

	response := []byte(fmt.Sprintf("Status: %v, Output: %v\n", status, output))
	log.Info().Msgf("response: %v", response)

	session.Publish(googleConf.ResponseTopicName, response)
}
