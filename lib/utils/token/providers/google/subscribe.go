package google

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (p *Provider) Get() string {
	log.Info().Msg("retrieving token from subscription")
	p.Conf.Subscribe(p.TokenTopicName, p.TokenSubscriptionName, p)

	log.Info().Msg("token retrieval complete")

	return p.GetToken()
}

func (p *Provider) MessageDeserialized(message []byte) {
	p.SetToken(string(message))

	log.Info().Msg("token processed successfully")
	logging.Warn(p.PublishStatus("token retrieved", true), "MessageDeserialized")

	p.Conf.Cancel()
}

func (p *Provider) MessageParseError(err error) {
	log.Error().Err(err).Msg("token retrieval failed")
	logging.Warn(p.PublishStatus("token retrieval error", false), "MessageParseError")
}
