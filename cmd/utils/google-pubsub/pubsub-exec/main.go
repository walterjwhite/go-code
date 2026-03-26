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
	Args []string
}

func init() {
	application.Configure(subscriberConf)
	if err := subscriberConf.PubSubConf.Init(application.Context); err != nil {
		logging.Error(fmt.Errorf("failed to initialize PubSub configuration: %v", err))
	}

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

	e := Executor{}

	subscriberConf.PubSubConf.Subscribe(subscriberConf.ExecutionTopicName, subscriberConf.ExecutionTopicName, &e)
	application.Wait()
}
