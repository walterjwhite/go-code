package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/google"

	"os"
)

type PublisherConfiguration struct {
	TopicName  string
	FileName   string `yaml:"file_name" flag:"file f" desc:"Publish the contents of a file instead of positional arguments"`
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

	payload, err := buildPayload()
	if err != nil {
		logging.Error(err)
		return
	}

	logging.Warn(publisherConfiguration.GoogleConf.Publish(publisherConfiguration.TopicName, payload), "main")
}

func buildPayload() ([]byte, error) {
	args := flag.Args()
	hasFile := len(publisherConfiguration.FileName) > 0
	hasArgs := len(args) > 0

	switch {
	case hasFile && hasArgs:
		return nil, errors.New("choose exactly one mode: provide --file or positional arguments")
	case hasFile:
		return readFilePayload(publisherConfiguration.FileName)
	case hasArgs:
		payload, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		return payload, nil
	default:
		return nil, errors.New("expected input: provide --file <path> or positional arguments")
	}
}

func readFilePayload(name string) ([]byte, error) {
	info, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, errors.New("file path must not be a directory")
	}

	return os.ReadFile(name)
}
