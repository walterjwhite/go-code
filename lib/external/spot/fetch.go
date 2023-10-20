package spot

import (
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/time/periodic"
)

func (c *Configuration) fetch() {
	// ensure we don't incorrectly set a refresh interval
	if c.RefreshInterval < minRefreshInterval {
		c.RefreshInterval = minRefreshInterval
	}

	c.fetchPeriodic = periodic.Periodic(c.ctx, &c.RefreshInterval, true, c.doFetch)
}

func (c *Configuration) doFetch() error {
	messages := c.feedFetcher.Fetch()

	newMessages := make([]*data.Record, 0)
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]

		log.Debug().Msgf("processing message: %v", message)

		if c.isMessageNew(message) {
			r := data.New(message)
			newMessages = append(newMessages, r)

			log.Info().Msgf("received new record: %v", r)

			c.writer.Write(r)
		}
	}

	for _, newMessage := range newMessages {
		c.onNewRecord(c.Session.LatestReceivedRecord, newMessage)

		c.Session.LatestReceivedRecord = newMessage
	}

	return nil
}

func (c *Configuration) isMessageNew(message *client.Message) bool {
	return c.Session.LatestReceivedRecord == nil || message.Id > c.Session.LatestReceivedRecord.Id
}
