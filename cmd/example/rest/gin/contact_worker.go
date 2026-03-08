package main

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"strings"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func contactWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	serviceURL := os.Getenv("PULSAR_URL")
	topic := os.Getenv("PULSAR_TOPIC")
	subscription := os.Getenv("PULSAR_SUBSCRIPTION")
	if subscription == "" {
		subscription = "contact-sub"
	}
	if serviceURL == "" || topic == "" {
		fmt.Println("PULSAR_URL and PULSAR_TOPIC must be set for the worker")
		os.Exit(1)
	}

	emailCfg, err := getEmailConfigFromEnv()
	if err != nil {
		fmt.Printf("failed to get email config: %v\n", err)
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               serviceURL,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		fmt.Printf("could not instantiate Pulsar client: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subscription,
		Type:             pulsar.Shared,
	})
	if err != nil {
		fmt.Printf("could not create consumer: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	fmt.Printf("Worker started: listening on topic=%s subscription=%s\n", topic, subscription)

	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			fmt.Printf("receive error: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var req ContactRequest
		if err := json.Unmarshal(msg.Payload(), &req); err != nil {
			fmt.Printf("invalid message payload: %v\n", err)
			if ackErr := consumer.Ack(msg); ackErr != nil {
				fmt.Printf("error acknowledging malformed message: %v\n", ackErr)
			}
			continue
		}

		if strings.TrimSpace(req.Email) == "" || !validateEmailAddress(req.Email) {
			fmt.Printf("invalid email in message, acknowledging and skipping\n")
			if ackErr := consumer.Ack(msg); ackErr != nil {
				fmt.Printf("error acknowledging - invalid email address in message: %v\n", ackErr)
			}
			continue
		}

		if err := sendContactEmail(emailCfg, req); err != nil {
			fmt.Printf("failed to send email: %v\n", err)
			consumer.Nack(msg)
			time.Sleep(1 * time.Second)
			continue
		}

		if ackErr := consumer.Ack(msg); ackErr != nil {
			fmt.Printf("error acknowledging successful delivery of message: %v\n", ackErr)
		}
		fmt.Printf("email sent and message acknowledged\n")
	}
}
