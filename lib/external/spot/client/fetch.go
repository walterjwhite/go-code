package client

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"net/http"
	"time"
)

var (
)

type FeedFetcher interface {
	Fetch() []*Message
}

func (f *Feed) Fetch() []*Message {
	url := fmt.Sprintf(*endpointUrlFlag, f.Id)

	f.initClient()

	resp, err := f.client.Get(url)
	logging.Panic(err)

	defer resp.Body.Close()

	container := &Container{}
	logging.Panic(json.NewDecoder(resp.Body).Decode(&container))

	log.Debug().Msgf("received: %v", container)

	return container.Response.FeedMessageResponse.Messages.Message
}

func (f *Feed) initClient() {
	if f.client == nil {
		f.client = &http.Client{
			Timeout: time.Duration(5 * time.Second),
		}
	}
}
