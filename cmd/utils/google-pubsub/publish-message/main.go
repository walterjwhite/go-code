package main

import (
	"errors"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/publisher/provider/google"

	"os"
)

var (
	provider *google.Provider
	initErr  error
)

func init() {
	provider, initErr = google.New(application.Context)
	if initErr != nil {
		logging.Error(initErr)
	}
}

func main() {
	if initErr != nil {
		logging.Error(errors.New("failed to initialize provider"))
		return
	}

	if len(os.Args) < 2 {
		logging.Error(errors.New("expected message"))
		return
	}

	const maxMessageSize = 10 * 1024 * 1024 // 10 MB limit
	message := os.Args[1]
	if len(message) > maxMessageSize {
		logging.Error(errors.New("message size exceeds maximum limit"))
		return
	}

	logging.Warn(provider.Publish([]byte(message)), "publish-event.Publish")
}
