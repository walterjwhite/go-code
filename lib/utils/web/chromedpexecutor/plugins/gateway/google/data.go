package google

type Provider struct {
	CredentialsFile string
	ProjectId       string

	TokenTopicName        string
	TokenSubscriptionName string

	StatusTopicName        string
	StatusSubscriptionName string
}
