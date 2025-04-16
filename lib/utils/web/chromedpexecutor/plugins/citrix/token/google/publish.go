package google

import (
	"fmt"
)

func (p *Provider) PublishToken(token string) {
	p.session.Publish(p.TokenTopicName, token)
}

func (p *Provider) publishStatus(status string, successful bool) {
	p.session.Publish(p.StatusTopicName, fmt.Sprintf("%s|%v", status, successful))
}
