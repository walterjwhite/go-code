package stdin

import (
	"context"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/wait"
	"io/ioutil"
	"os"
	"time"
)

/* Attempts to read the specified file, periodically polling <Interval> and waiting at most <Timeout>*/
type FileReader struct {
	Filename string
	Timeout  *time.Duration
	Interval *time.Duration
	Context  context.Context
}

func (r *FileReader) Get() string {
	wait.Wait(r.Context, r.Interval, r.Timeout, r.fileExists)
	return r.read()
}

func (r *FileReader) fileExists() bool {
	info, err := os.Stat(r.Filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func (r *FileReader) read() string {
	data, err := ioutil.ReadFile(r.Filename)
	logging.Panic(err)

	return string(data)
}
