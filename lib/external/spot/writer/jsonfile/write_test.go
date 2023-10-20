package jsonfile

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	tmpDirFile, err := os.CreateTemp("", "test_export")
	logging.Panic(err)

	defer os.RemoveAll(tmpDirFile.Name())

	c := &data.Session{FeedId: "export_test"}
	c.SessionPath = tmpDirFile.Name()

	r := &data.Record{Id: 1, Latitude: 45.0, Longitude: -80.0, Message: "This is a test", MessageType: client.OK}

	a := &RecordAppenderConfiguration{Session: c}
	a.Write(r)
}
