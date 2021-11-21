package gpx

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"

	"io/ioutil"
	"os"
	"path/filepath"
)

type LatestFile struct {
	Path     string
	FileInfo os.FileInfo
}

func Latest(s *data.Session) []*data.Record {
	f := &LatestFile{}

	latest(s.DataPath, f)
	return get(filepath.Join(f.Path, f.FileInfo.Name()))
}

func latest(path string, f *LatestFile) {
	log.Debug().Msgf("in: %s", path)
	files, err := ioutil.ReadDir(path)
	logging.Panic(err)

	for _, file := range files {
		log.Debug().Msgf("inspecting: %s", file.Name())

		if file.IsDir() {
			latest(filepath.Join(path, file.Name()), f)
		} else {
			if f.FileInfo == nil {
				f.Path = path
				f.FileInfo = file

				continue
			}

			if file.ModTime().After(f.FileInfo.ModTime()) {
				f.Path = path
				f.FileInfo = file

				continue
			}
		}
	}
}
