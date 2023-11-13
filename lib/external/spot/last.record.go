package spot

import (
	"bufio"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"os"

	"path/filepath"
	"time"
)

func (c *Configuration) initLastRecord() {
	if _, err := os.Stat(c.Session.SessionPath); os.IsNotExist(err) {
		log.Warn().Msgf("session path !exist: %v", c.Session.SessionPath)
		return
	}

	if _, err := os.Stat(c.Session.DataPath); os.IsNotExist(err) {
		log.Warn().Msgf("session path !exist: %v", c.Session.DataPath)
		return
	}

	l := &lastRecord{}
	c.getLastRecordFilename(c.Session.DataPath, l)

	if len(l.Filename) == 0 {
		log.Warn().Msgf("Existing session data not found")
		return
	}

	log.Debug().Msgf("last record filename: %s", l.Filename)
	file, err := os.Open(l.Filename)
	logging.Panic(err)

	defer file.Close()

	c.Session.LatestReceivedRecord = &data.Record{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logging.Panic(json.Unmarshal([]byte(scanner.Text()), c.Session.LatestReceivedRecord))
		logging.Panic(scanner.Err())

		log.Debug().Msgf("last record: %v", c.Session.LatestReceivedRecord)
	}
}

type lastRecord struct {
	Filename   string
	ModifyTime time.Time
}

func (c *Configuration) getLastRecordFilename(path string, l *lastRecord) {
	files, err := os.ReadDir(path)
	logging.Panic(err)

	for _, file := range files {
		fileInfo, err := file.Info()
		logging.Panic(err)

		if fileInfo.Mode().IsRegular() {
			if isNewer(l, fileInfo) {
				l.ModifyTime = fileInfo.ModTime()
				l.Filename = filepath.Join(path, file.Name())
			}
		} else {
			c.getLastRecordFilename(filepath.Join(path, file.Name()), l)
		}
	}
}

func isNewer(l *lastRecord, file os.FileInfo) bool {
	if len(l.Filename) == 0 {
		return true
	}

	return !file.ModTime().Before(l.ModifyTime) &&
		file.ModTime().After(l.ModifyTime)
}
