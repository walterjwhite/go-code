package main

import (
	"errors"

	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/exec"
	"github.com/walterjwhite/go-code/lib/net/google"
	"github.com/walterjwhite/go-code/lib/security/encryption/aes"
	"github.com/walterjwhite/go-code/lib/security/encryption/providers/file"

	oexec "os/exec"
)

type SubscriberConfiguration struct {
	CredentialsFile string
	ProjectId       string

	TopicName        string
	SubscriptionName string

	ResponseTopicName     string
	EncryptionKeyFilename string
}

var (
	googleConf = &SubscriberConfiguration{}
	aesConf    = &aes.Configuration{}
	session    *google.Session

	cmd = flag.String("cmd", "", "cmd to execute on receipt of pubsub message")
)

type Executor struct {
	Cmd *exec.Cmd
}

func init() {
	application.Configure(googleConf)

	aesConf.Encryption = file.New(googleConf.EncryptionKeyFilename)
}

func main() {
	if len(*cmd) == 0 {
		logging.Panic(errors.New("-cmd=<COMMAND>"))
	}

	e := Executor{}

	session = google.New(googleConf.CredentialsFile, googleConf.ProjectId, application.Context)
	session.AesConf = aesConf
	session.EnableCompression = true

	session.Subscribe(googleConf.TopicName, googleConf.SubscriptionName, &e)
	application.Wait()
}

func (e *Executor) New() any {
	e.Cmd = &exec.Cmd{}
	return e.Cmd
}

func (e *Executor) MessageDeserialized() {
	log.Info().Msgf("running: %s -> %s", e.Cmd.FunctionName, e.Cmd.Args)

	e.Cmd.Args = append([]string{e.Cmd.FunctionName}, e.Cmd.Args...)

	ecmd := oexec.Command(*cmd, e.Cmd.Args...)
	output, err := ecmd.Output()

	status := 0
	if err != nil {
		if exitError, ok := err.(*oexec.ExitError); ok {
			status = exitError.ExitCode()
			log.Warn().Msgf("Error running: %s (%s) -> %v", *cmd, e.Cmd.Args, status)
		}
	}

	respond(status, string(output))
}

func (e *Executor) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}

func respond(status int, output string) {
	log.Info().Msgf("status: %v, output: %s", status, output)

	response := fmt.Sprintf("Status: %v, Output: \n%v\n", status, output)
	log.Info().Msgf("response: %v", response)

	session.Publish(googleConf.ResponseTopicName, response)
}
