package gpx

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	"io/ioutil"

	"path/filepath"
	"sort"
)

func All(s *data.Session) []*data.Record {
	return dir(s.DataPath)
}

func dir(path string) []*data.Record {
	records := make([]*data.Record, 0)

	files, err := ioutil.ReadDir(path)
	logging.Panic(err)

	// sort files
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for _, file := range files {
		filename := filepath.Join(path, file.Name())

		if file.IsDir() {
			records = append(records, dir(filename)...)
		} else {
			records = append(records, get(filename)...)
		}
	}

	return records
}
