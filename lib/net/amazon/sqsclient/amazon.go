package sqsclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Client struct {
	svc      *sqs.Client
	queueURL string
}

func New(ctx context.Context, region, queueURL string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("sqsclient: load AWS config: %w", err)
	}
	return &Client{
		svc:      sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func NewWithConfig(cfg aws.Config, queueURL string) *Client {
	return &Client{
		svc:      sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}
}

func (c *Client) ReceiveMessage(
	ctx context.Context,
	visibilityTimeout, waitSeconds int32,
) ([]byte, error) {
	if visibilityTimeout == 0 {
		visibilityTimeout = 30
	}
	if waitSeconds == 0 {
		waitSeconds = 20
	}

	out, err := c.svc.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   visibilityTimeout,
		WaitTimeSeconds:     waitSeconds,
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("sqsclient: receive message: %w", err)
	}

	if len(out.Messages) == 0 {
		return nil, nil // queue empty within the wait window
	}

	msg := out.Messages[0]

	_, err = c.svc.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		return nil, fmt.Errorf("sqsclient: delete message: %w", err)
	}

	if msg.Body == nil {
		return []byte{}, nil
	}
	return []byte(*msg.Body), nil
}

func (c *Client) PublishMessage(ctx context.Context, payload []byte) (string, error) {
	out, err := c.svc.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(payload)),
	})
	if err != nil {
		return "", fmt.Errorf("sqsclient: send message: %w", err)
	}
	if out.MessageId == nil {
		return "", fmt.Errorf("sqsclient: send message: empty message ID in response")
	}
	return *out.MessageId, nil
}
