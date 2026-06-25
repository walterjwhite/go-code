package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/ai/rag"
	"github.com/walterjwhite/go-code/lib/ai/rag/pubsub"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/google"
)

type Conf struct {
	RagConfig           *rag.Config
	PubsubConfiguration *google.Conf

	RequestTopic        string
	RequestSubscription string
	ResponseTopic       string

	QueryTimeout time.Duration

	service *pubsub.Service
}

var (
	config = &Conf{}
)

func init() {
	application.Configure(config)
	validate()

	logging.Error(config.PubsubConfiguration.Init(application.Context), "initialise pubsub")

	var err error
	log.Info().Msg("connecting to Qdrant, Ollama, and Pub/Sub")
	logging.Error(config.RagConfig.New(application.Context), "initialize engine")

	log.Info().Msgf("rag configuration: %+v", *config.RagConfig)

	config.service, err = pubsub.New(pubsub.Config{
		Context:             application.Context,
		PubSub:              config.PubsubConfiguration,
		RequestTopic:        config.RequestTopic,
		RequestSubscription: config.RequestSubscription,
		ResponseTopic:       config.ResponseTopic,
		QueryTimeout:        config.QueryTimeout,
	}, config.RagConfig.Engine)
	logging.Error(err, "initialise pubsub rag service")

	log.Info().Msgf("pubsub configuration: %+v", *config.PubsubConfiguration)
}

func main() {
	defer application.OnPanic()


	defer config.PubsubConfiguration.Cancel()

	log.Info().Msgf("listening on subscription: %s", config.RequestSubscription)
	log.Info().Msgf("publishing responses to topic: %s", config.ResponseTopic)

	config.service.Run()
}

func validate() {
	if config.RequestSubscription == "" {
		log.Fatal().Msg("request subscription is required")
	}
	if config.ResponseTopic == "" {
		log.Fatal().Msg("response topic is required")
	}
	if config.PubsubConfiguration.ProjectId == "" {
		log.Fatal().Msg("gcp project ID is required")
	}
	if config.PubsubConfiguration.CredentialsFile == "" {
		log.Fatal().Msg("gcp credentials file is required")
	}
}
