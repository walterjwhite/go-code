package main

import (
	"fmt"
	"os"
	"strconv"
)

func getEmailConfigFromEnv() (*EmailConfig, error) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	to := os.Getenv("CONTACT_TO")

	if host == "" || portStr == "" {
		return nil, fmt.Errorf("SMTP_HOST and SMTP_PORT must be set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	if from == "" {
		from = username
	}
	if to == "" {
		to = os.Getenv("USER")
		if to == "" {
			to = "contact@example.com"
		}
	}

	return &EmailConfig{
		SMTPHost: host,
		SMTPPort: port,
		Username: username,
		Password: password,
		From:     from,
		To:       to,
	}, nil
}

func getPulsarConfigFromEnv() (serviceURL, topic, subscription string, err error) {
	serviceURL = os.Getenv("PULSAR_URL")
	topic = os.Getenv("PULSAR_TOPIC")
	subscription = os.Getenv("PULSAR_SUBSCRIPTION") // not required for producer but keep for docs

	if serviceURL == "" {
		err = fmt.Errorf("PULSAR_URL must be set")
		return
	}
	if topic == "" {
		err = fmt.Errorf("PULSAR_TOPIC must be set")
		return
	}
	if subscription == "" {
		err = fmt.Errorf("PULSAR_SUBSCRIPTION must be set")
		return
	}
	return
}
