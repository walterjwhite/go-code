package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	DatabaseURL string
	AppPort     string
	LogLevel    zerolog.Level
	ReadTimeout time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	level := zerolog.InfoLevel
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if l, err := zerolog.ParseLevel(lvl); err == nil {
			level = l
		}
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://$USERNAME:$PASSWORD@$HOSTNAME:$PORT/$DATABASE?sslmode=disable"
	}

	return &Config{
		DatabaseURL: dbURL,
		AppPort:     port,
		LogLevel:    level,
		ReadTimeout: 10 * time.Second,
	}, nil
}
