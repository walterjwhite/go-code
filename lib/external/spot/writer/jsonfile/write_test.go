package jsonfile

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/client"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"io/ioutil"
	"os"
	"testing"
)

func TestExport(t *testing.T) {
	tmpDirName, err := ioutil.TempDir("", "test_export")
	logging.Panic(err)

	defer os.RemoveAll(tmpDirName)

	c := &data.Session{FeedId: "export_test"}
	c.SessionPath = tmpDirName

	r := &data.Record{Id: 1, Latitude: 45.0, Longitude: -80.0, Message: "This is a test", MessageType: client.OK}

	a := &RecordAppenderConfiguration{Session: c}
	a.Write(r)
}
