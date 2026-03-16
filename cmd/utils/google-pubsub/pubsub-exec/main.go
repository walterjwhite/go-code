package main

import (
	"errors"

	"flag"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"
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
	Args []string
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
