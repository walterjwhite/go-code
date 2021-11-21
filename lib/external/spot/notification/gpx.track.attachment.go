package notification

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/external/spot/gpx"

	"github.com/walterjwhite/go-code/lib/net/email"
	"io/ioutil"
	"os"
)

func (c *Notification) getTrackAsAttachment() *email.EmailAttachment {
	return &email.EmailAttachment{Name: fmt.Sprintf("latest-%s-track.gpx", c.Session.FeedId), Data: bytes.NewBuffer(c.exportLatestTrack())}
}

func (c *Notification) exportLatestTrack() []byte {
	tmpFile, err := ioutil.TempFile(os.TempDir(), c.Session.FeedId)
	logging.Panic(err)

	defer os.Remove(tmpFile.Name())

	// export the latest data
	records := gpx.Latest(c.Session)
	// latestRecord := records[len(records)-1]

	exportFilename := gpx.Export(records, tmpFile.Name())
	defer os.Remove(exportFilename)

	// read data to []byte
	data, err := ioutil.ReadFile(exportFilename)
	logging.Panic(err)

	log.Debug().Msgf("attachment size: %v", len(data))

	return data
}
