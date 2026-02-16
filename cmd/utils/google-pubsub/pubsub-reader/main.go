package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"

	"github.com/walterjwhite/go-code/lib/net/google"
)

type ReadSubscriberConfiguration struct {
	TopicName        string
	SubscriptionName string
	PubSubConf       *google.Conf
}

type Callback struct {
}

var (
	readSubscriberConf = &ReadSubscriberConfiguration{}
)

func init() {
	application.Configure(readSubscriberConf)
	readSubscriberConf.PubSubConf.Init(application.Context)
}

func main() {
	defer application.OnPanic()
	c := &Callback{}

	readSubscriberConf.PubSubConf.Subscribe(readSubscriberConf.TopicName, readSubscriberConf.SubscriptionName, c)
	application.Wait()
}

func (c *Callback) MessageDeserialized(deserialized []byte) {
	log.Info().Msgf("callback: %s", string(deserialized))
}

func (c *Callback) MessageParseError(err error) {
	log.Error().Msgf("Error parsing message: %v", err)
}
