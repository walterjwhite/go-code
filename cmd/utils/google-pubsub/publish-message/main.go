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
)

func init() {
	provider = google.New(application.Context)
}

func main() {
	if len(os.Args) < 2 {
		logging.Error(errors.New("expected message"))
		return
	}

	logging.Warn(provider.Publish([]byte(os.Args[1])), "publish-event.Publish")
}
