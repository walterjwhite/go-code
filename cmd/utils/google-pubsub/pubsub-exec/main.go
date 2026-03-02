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
	"regexp"
	"unicode"
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
	if err := subscriberConf.PubSubConf.Init(application.Context); err != nil {
		logging.Error(fmt.Errorf("failed to initialize PubSub configuration: %v", err))
	}
}

func main() {
	defer application.OnPanic()
	if len(*cmd) == 0 {
		logging.Error(errors.New("-cmd=<COMMAND>"))
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

	if !isValidCommandName(e.Cmd.FunctionName) {
		log.Warn().Msgf("invalid function name: %s", e.Cmd.FunctionName)
		return
	}

	for _, arg := range e.Cmd.Args {
		if !isValidArgument(arg) {
			log.Warn().Msgf("invalid argument detected: %s", arg)
			return
		}
	}

	e.Cmd.Args = append([]string{e.Cmd.FunctionName}, e.Cmd.Args...)

	ecmd := oexec.Command(*cmd, e.Cmd.Args...) // #nosec
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

func isValidCommandName(name string) bool {
	if len(name) == 0 || len(name) > 256 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_\-.]+$`, name)
	return matched
}

func isValidArgument(arg string) bool {
	const maxArgLength = 4096
	if len(arg) > maxArgLength {
		return false
	}

	for _, r := range arg {
		if !isValidCharacter(r) {
			return false
		}
	}
	return true
}

func isValidCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) ||
		unicode.IsSpace(r) ||
		r == '.' || r == ',' || r == '-' || r == '_' ||
		r == '+' || r == '='
}

func (e *Executor) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}

func respond(status int, output string) {
	log.Info().Msgf("status: %v", status)

	response := fmt.Sprintf("Status: %v, Output: \n%v\n", status, output)
	log.Debug().Msgf("response published with status: %v", status)

	logging.Warn(subscriberConf.PubSubConf.Publish(subscriberConf.ResponseTopicName, []byte(response)), "respond")
}
