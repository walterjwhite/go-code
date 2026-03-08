package main

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	_ = os.Unsetenv("APP_PORT")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("DATABASE_URL")

	cfg, err := Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
	assert.Contains(t, cfg.DatabaseURL, "postgres://")
}

func TestLoad_CustomPort(t *testing.T) {
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Setenv("APP_PORT", "3000")
	defer func() { _ = os.Unsetenv("APP_PORT") }()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "3000", cfg.AppPort)
}

func TestLoad_CustomLogLevel(t *testing.T) {
	_ = os.Unsetenv("APP_PORT")
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Setenv("LOG_LEVEL", "debug")
	defer func() { _ = os.Unsetenv("LOG_LEVEL") }()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, zerolog.DebugLevel, cfg.LogLevel)
}

func TestLoad_CustomLogLevelError(t *testing.T) {
	_ = os.Unsetenv("APP_PORT")
	_ = os.Unsetenv("DATABASE_URL")
	_ = os.Setenv("LOG_LEVEL", "invalid")
	defer func() { _ = os.Unsetenv("LOG_LEVEL") }()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel) // defaults to info on invalid
}

func TestLoad_CustomDatabaseURL(t *testing.T) {
	_ = os.Unsetenv("APP_PORT")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	defer func() { _ = os.Unsetenv("DATABASE_URL") }()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", cfg.DatabaseURL)
}

func TestLoad_AllCustomValues(t *testing.T) {
	_ = os.Setenv("APP_PORT", "9000")
	_ = os.Setenv("LOG_LEVEL", "warn")
	_ = os.Setenv("DATABASE_URL", "sqlite://test.db")
	defer func() {
		_ = os.Unsetenv("APP_PORT")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("DATABASE_URL")
	}()

	cfg, err := Load()

	assert.NoError(t, err)
	assert.Equal(t, "9000", cfg.AppPort)
	assert.Equal(t, zerolog.WarnLevel, cfg.LogLevel)
	assert.Equal(t, "sqlite://test.db", cfg.DatabaseURL)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
}

func TestConfig_Struct(t *testing.T) {
	cfg := &Config{
		DatabaseURL: "postgres://localhost",
		AppPort:     "8080",
		LogLevel:    zerolog.InfoLevel,
		ReadTimeout: 5 * time.Second,
	}

	assert.Equal(t, "postgres://localhost", cfg.DatabaseURL)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, zerolog.InfoLevel, cfg.LogLevel)
	assert.Equal(t, 5*time.Second, cfg.ReadTimeout)
}
