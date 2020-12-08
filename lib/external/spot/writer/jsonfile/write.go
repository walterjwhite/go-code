package jsonfile

import (
	"encoding/json"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/external/spot/data"
	"github.com/walterjwhite/go/lib/time/timeformatter/day"
	"os"
	"path"
	"path/filepath"
)

func (a *RecordAppenderConfiguration) Write(r *data.Record) {
	// TODO: use the date/time from the record itself
	filename := path.Join(a.Session.DataPath, day.Get())
	logging.Panic(os.MkdirAll(filepath.Dir(filename), os.ModePerm))

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	logging.Panic(err)

	defer f.Close()

	// serialize to json
	data, err := json.Marshal(r)
	logging.Panic(err)

	_, err = f.WriteString(string(data) + "\n")
	logging.Panic(err)
}
