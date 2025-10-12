package google

import (
	"fmt"
)

func (p *Provider) PublishToken(token string) error {
	return p.Conf.Publish(p.TokenTopicName, []byte(token))
}

func (p *Provider) PublishStatus(status string, successful bool) error {
	return p.Conf.Publish(p.StatusTopicName, []byte(fmt.Sprintf("%s|%v", status, successful)))
}
