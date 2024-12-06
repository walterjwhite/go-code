package main

import (
	"flag"
	"fmt"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"github.com/walterjwhite/go-code/lib/external/spot/gpx"
	"github.com/walterjwhite/go-code/lib/time/timeformatter/timestamp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

func export(c *spot.Configuration) {
	var records []*data.Record

	var action string
	if isExportLatest() {
		records = gpx.Latest(c.Session)
		action = "latest"
	} else if strings.Compare(flag.Args()[0], "all") == 0 {
		records = gpx.All(c.Session)
		action = "all"
	} else {
		date, err := time.Parse(dateFormat, flag.Args()[0])
		logging.Panic(err)

		records = gpx.Day(c.Session, &date)
		action = flag.Args()[0]
	}

	gpx.Export(records, getFilename(action))
}

func isExportLatest() bool {
	if len(flag.Args()) == 0 {
		return true
	}

	return strings.Compare(flag.Args()[0], "latest") == 0
}

func getFilename(action string) string {
	wd, err := os.Getwd()
	logging.Panic(err)

	return filepath.Join(wd, fmt.Sprintf("%v-%v", action, timestamp.Get()))
}
