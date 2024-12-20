package notification

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/external/spot/gpx"

	"github.com/walterjwhite/go-code/lib/net/email"

	"os"
)

func (c *Notification) getTrackAsAttachment() *email.EmailAttachment {
	return &email.EmailAttachment{Name: fmt.Sprintf("latest-%s-track.gpx", c.Session.FeedId), Data: bytes.NewBuffer(c.exportLatestTrack())}
}

func (c *Notification) exportLatestTrack() []byte {
	tmpFile, err := os.CreateTemp(os.TempDir(), c.Session.FeedId)
	logging.Panic(err)

	defer os.Remove(tmpFile.Name())

	records := gpx.Latest(c.Session)

	exportFilename := gpx.Export(records, tmpFile.Name())
	defer os.Remove(exportFilename)

	data, err := os.ReadFile(exportFilename)
	logging.Panic(err)

	log.Debug().Msgf("attachment size: %v", len(data))

	return data
}
