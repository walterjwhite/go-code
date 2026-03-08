package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"time"
)

func publishContactMessageToPulsar(req ContactRequest) error {
	serviceURL, topic, _, err := getPulsarConfigFromEnv()
	if err != nil {
		return fmt.Errorf("failed to get pulsar config: %w", err)
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               serviceURL,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to create pulsar client: %w", err)
	}

	defer client.Close()

	prod, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer prod.Close()

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := &pulsar.ProducerMessage{
		Payload: payload,
		Key:     req.Email,
		Properties: map[string]string{
			"source": "contact-form",
		},
	}

	_, err = prod.Send(ctx, msg)
	return err
}
