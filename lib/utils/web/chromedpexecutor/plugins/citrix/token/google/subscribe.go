package google

import (
	"fmt"
	"github.com/rs/zerolog/log"
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
	p.PublishStatus(fmt.Sprintf("read token: %s", p.token), true)

	p.Conf.Cancel()
}

func (p *Provider) MessageParseError(err error) {
	log.Error().Msgf("error reading token: %s", p.token)
	p.PublishStatus(fmt.Sprintf("error reading message: %s", err), false)
}
