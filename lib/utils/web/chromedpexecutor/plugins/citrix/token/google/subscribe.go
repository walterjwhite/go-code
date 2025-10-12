package google

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func (p *Provider) ReadToken() *string {
	log.Info().Msg("reading token")
	p.Conf.Subscribe(p.TokenTopicName, p.TokenSubscriptionName, p)

	log.Info().Msgf("read token: %s", p.token)

	return &p.token
}

func (p *Provider) MessageDeserialized(message []byte) {
	p.token = string(message)

	log.Info().Msgf("read token: %s", p.token)
	logging.Warn(p.PublishStatus(fmt.Sprintf("read token: %s", p.token), true), false, "MessageDeserialized")

	p.Conf.Cancel()
}

func (p *Provider) MessageParseError(err error) {
	log.Error().Msgf("error reading token: %s", p.token)
	logging.Warn(p.PublishStatus(fmt.Sprintf("error reading message: %s", err), false), false, "MessageParseError")
}
