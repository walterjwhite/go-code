package gpx

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"
	

	"path/filepath"
	"sort"
	"os"
)

func All(s *data.Session) []*data.Record {
	return dir(s.DataPath)
}

func dir(path string) []*data.Record {
	records := make([]*data.Record, 0)

	files, err := os.ReadDir(path)
	logging.Panic(err)

	// sort files
	sort.Slice(files, func(i, j int) bool {
		leftFileInfo, err := files[i].Info();
		logging.Panic(err)
		rightFileInfo, err := files[j].Info();
		logging.Panic(err)

		return leftFileInfo.ModTime().Before(rightFileInfo.ModTime())
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
