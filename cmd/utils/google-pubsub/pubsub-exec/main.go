package main

import (
	"encoding/json"
	"errors"

	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/exec"
	"github.com/walterjwhite/go-code/lib/net/google"

	oexec "os/exec"
)

type SubscriberConfiguration struct {
	TopicName        string
	SubscriptionName string

	ResponseTopicName string
	PubSubConf        *google.Conf
}

var (
	subscriberConf = &SubscriberConfiguration{}

	cmd = flag.String("cmd", "", "cmd to execute on receipt of pubsub message")
)

type Executor struct {
	Cmd *exec.Cmd
}

func init() {
	application.Configure(subscriberConf)
	subscriberConf.PubSubConf.Init(application.Context)
}

func main() {
	if len(*cmd) == 0 {
		logging.Panic(errors.New("-cmd=<COMMAND>"))
	}

	e := Executor{}

	subscriberConf.PubSubConf.Subscribe(subscriberConf.TopicName, subscriberConf.SubscriptionName, &e)
	application.Wait()
}

func (e *Executor) MessageDeserialized(deserialized []byte) {
	e.Cmd = &exec.Cmd{}

	err := json.Unmarshal(deserialized, e.Cmd)
	if err != nil {
		log.Warn().Msgf("error converting to exec.cmd, %v", err)
		return
	}

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
	} else {
		log.Info().Msgf("Successfully ran: %s (%s) -> %v", *cmd, e.Cmd.Args, status)
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

	logging.Warn(subscriberConf.PubSubConf.Publish(subscriberConf.ResponseTopicName, []byte(response)), "respond")
}
