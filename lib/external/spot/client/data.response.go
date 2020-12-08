package client

type Container struct {
	Response *Response
}

type Response struct {
	FeedMessageResponse *FeedMessageResponse
}

type FeedMessageResponse struct {
	Count         int
	Feed          *Feed
	TotalCount    int
	ActivityCount int
	Messages      *MessageWrapper
}

type MessageWrapper struct {
	Message []*Message
}
