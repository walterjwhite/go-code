package main

import (
	"encoding/json"
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"

	"os"
)

type PublisherConfiguration struct {
	TopicName  string
	GoogleConf *google.Conf
}

var (
	publisherConfiguration = &PublisherConfiguration{}
)

func init() {
	application.Configure(publisherConfiguration)
	if err := publisherConfiguration.GoogleConf.Init(application.Context); err != nil {
		logging.Error(err)
	}
}

func main() {
	defer application.OnPanic()

	if len(os.Args) == 1 {
		logging.Error(errors.New("expecting arguments, at least a command/function must be provided"))
	}

	jsonString, err := json.Marshal(os.Args[1:])
	if err != nil {
		logging.Warn(err, "unable to convert to json")
		return
	}

	logging.Warn(publisherConfiguration.GoogleConf.Publish(publisherConfiguration.TopicName, jsonString), "main")
}
