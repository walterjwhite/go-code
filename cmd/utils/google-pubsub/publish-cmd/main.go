package main

import (
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

	functionName = flag.String("functionName", "", "function to execute remotely")
	arguments    = flag.String("arguments", "", "arguments to pass functionName, optional")
)

func init() {
	application.Configure(publisherConfiguration)
	publisherConfiguration.GoogleConf.Init(application.Context)
}

func main() {
	if len(*functionName) == 0 {
		logging.Panic(errors.New("expecting command to be non-empty"))
	}

	c := exec.Cmd{FunctionName: *functionName}
	if len(*arguments) != 0 {
		c.Args = strings.Fields(*arguments)
	}

	logging.Warn(publisherConfiguration.GoogleConf.Publish(publisherConfiguration.TopicName, []byte("TODO: convert c to string using json, then return byte array")), false, "main")
}
