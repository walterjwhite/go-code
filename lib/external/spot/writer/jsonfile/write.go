package jsonfile

import (
	"encoding/json"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/time/timeformatter/day"
	"os"
	"path"
	"path/filepath"
)

func (a *RecordAppenderConfiguration) Write(r *data.Record) {
	filename := path.Join(a.Session.DataPath, day.Get())
	logging.Panic(os.MkdirAll(filepath.Dir(filename), os.ModePerm))

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	logging.Panic(err)

	defer f.Close()

	data, err := json.Marshal(r)
	logging.Panic(err)

	_, err = f.WriteString(string(data) + "\n")
	logging.Panic(err)
}
