package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application/logging"
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
	logging.Error(err)

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
	logging.Error(err)
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
			logging.Warn(consumer.Ack(msg), "error acking malformed message")
			continue
		}

		if strings.TrimSpace(req.Email) == "" || !validateEmailAddress(req.Email) {
			fmt.Printf("invalid email in message, acking and skipping: %v\n", req.Email)
			logging.Warn(consumer.Ack(msg), "error acking - invalid email address in message")
			continue
		}

		if err := sendContactEmail(emailCfg, req); err != nil {
			fmt.Printf("failed to send email for %s: %v\n", req.Email, err)
			consumer.Nack(msg)
			time.Sleep(1 * time.Second)
			continue
		}

		logging.Warn(consumer.Ack(msg), "error acking successful delivery of message")
		fmt.Printf("email sent and message acked for %s\n", req.Email)
	}
}
