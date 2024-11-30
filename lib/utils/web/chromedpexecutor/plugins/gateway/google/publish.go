package google

import (
	"fmt"
	google_api "github.com/walterjwhite/go-code/lib/net/google"
)

func (p *Provider) PublishToken(token string) {
	google_api.Publish(p.CredentialsFile, p.ProjectId, p.TokenTopicName, token)
}

func (p *Provider) PublishStatus(status string, successful bool) {
	google_api.Publish(p.CredentialsFile, p.ProjectId, p.StatusTopicName, fmt.Sprintf("%s|%v", status, successful))
}
