package main

import (
	"errors"

	"flag"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"

	"os"
)

type SubscriberConfiguration struct {
	ExecutionTopicName string
	StatusTopicName    string

	PubSubConf *google.Conf
}

var (
	subscriberConf = &SubscriberConfiguration{}

	cmd = flag.String("cmd", "", "cmd to execute on receipt of pubsub message")
)

type Executor struct {
	Args    []string
	Runtime *PubSubRuntime
}

func init() {
	application.Configure(subscriberConf)

	name, err := os.Hostname()
	logging.Error(err, "hostname")

	subscriberConf.ExecutionTopicName = name + "_exec"
	subscriberConf.StatusTopicName = name + "_status"
}

func main() {
	defer application.OnPanic()
	if len(*cmd) == 0 {
		logging.Error(errors.New("-cmd=<COMMAND>"))
	}

	runtime, err := NewPubSubRuntime(application.Context, subscriberConf.PubSubConf)
	if err != nil {
		logging.Error(fmt.Errorf("failed to initialize PubSub runtime: %v", err))
	}

	e := Executor{Runtime: runtime}

	runtime.Subscribe(subscriberConf.ExecutionTopicName, &e)
	application.Wait()
}
