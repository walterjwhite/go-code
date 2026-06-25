package chat

const DefaultResponseBuffer = 64

func validateConfig(config Config) error {
	if config.PublishTopic == "" {
		return ErrPublishTopicRequired
	}
	if config.ResponseTopic == "" {
		return ErrResponseTopicRequired
	}
	if config.ResponseSubscription == "" {
		return ErrResponseSubRequired
	}
	if config.RateLimit.Messages <= 0 {
		return ErrRateLimitCountInvalid
	}
	if config.RateLimit.Window <= 0 {
		return ErrRateLimitWindowInvalid
	}
	return nil
}

func responseBufferSize(config Config) int {
	if config.ResponseBuffer > 0 {
		return config.ResponseBuffer
	}
	return DefaultResponseBuffer
}
