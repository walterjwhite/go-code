package main

import (
	"encoding/json"
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/net/exec"
	"github.com/walterjwhite/go-code/lib/net/google"

	"flag"
	"strings"
)

type PublisherConfiguration struct {
	TopicName  string
	GoogleConf *google.Conf
}

var (
	publisherConfiguration = &PublisherConfiguration{}

	functionName = flag.String("func", "", "function to execute")
	arguments    = flag.String("args", "", "arguments to pass function, optional")

	allowedFunctions = map[string]bool{
		"process":   true,
		"publish":   true,
		"subscribe": true,
	}
)

func init() {
	application.Configure(publisherConfiguration)
	if err := publisherConfiguration.GoogleConf.Init(application.Context); err != nil {
		logging.Error(err)
	}
}

func main() {
	defer application.OnPanic()
	flag.Parse()

	if len(*functionName) == 0 {
		logging.Error(errors.New("expecting command to be non-empty"))
	}

	if !isAllowedFunction(*functionName) {
		logging.Error(errors.New("function not in allowed list"))
	}

	c := exec.Cmd{FunctionName: *functionName}
	if len(*arguments) != 0 {
		c.Args = strings.Fields(*arguments)
	}

	jsonString, err := json.Marshal(c)
	if err != nil {
		logging.Warn(err, "unable to convert to json")
		return
	}

	logging.Warn(publisherConfiguration.GoogleConf.Publish(publisherConfiguration.TopicName, jsonString), "main")
}

func isAllowedFunction(functionName string) bool {
	return allowedFunctions[functionName]
}
