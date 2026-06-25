package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/ai/rag/pubsub"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/google"
)

type Conf struct {
	GoogleConf *google.Conf
	Topic      string
}

var (
	conf = &Conf{}
)

func init() {
	application.Configure(conf)
	logging.Error(conf.GoogleConf.Init(application.Context), "init")
}

func main() {
	defer application.OnPanic()

	if len(os.Args) < 2 {
		fmt.Println("Usage: query <question>")
		return
	}

	question := strings.Join(os.Args[1:], " ")

	m, err := json.Marshal(pubsub.Message{IPAddress: "cli", Message: question})
	logging.Error(err, "json.Marshal")

	logging.Warn(conf.GoogleConf.Publish(conf.Topic, m), "main")

	log.Info().Msgf("published message: %v to topic: %s", m, conf.Topic)
}
