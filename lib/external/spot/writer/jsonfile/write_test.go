package jsonfile

import (
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"os"
	"testing"
)

func TestExport(t *testing.T) {
	tmpDirFile, err := os.MkdirTemp("", "test_export")
	logging.Panic(err)

	defer os.RemoveAll(tmpDirFile)

	c := &data.Session{FeedId: "export_test"}
	c.DataPath = tmpDirFile

	r := &data.Record{Id: 1, Latitude: 45.0, Longitude: -80.0, Message: "This is a test", MessageType: client.OK}

	a := &RecordAppenderConfiguration{Session: c}
	log.Warn().Msgf("logging to %v", c.DataPath)

	a.Write(r)
}
