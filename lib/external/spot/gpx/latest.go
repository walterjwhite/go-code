package gpx

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/external/spot/data"

	
	"os"
	"path/filepath"
)

type LatestFile struct {
	Path     string
	DirEntry os.DirEntry
}

func Latest(s *data.Session) []*data.Record {
	f := &LatestFile{}

	latest(s.DataPath, f)
	return get(filepath.Join(f.Path, f.DirEntry.Name()))
}

func latest(path string, f *LatestFile) {
	log.Debug().Msgf("in: %s", path)
	files, err := os.ReadDir(path)
	logging.Panic(err)

	for _, file := range files {
		log.Debug().Msgf("inspecting: %s", file.Name())

		if file.IsDir() {
			latest(filepath.Join(path, file.Name()), f)
		} else {
			if f.DirEntry == nil {
				f.Path = path
				f.DirEntry = file

				continue
			}

			leftFileInfo, err := file.Info();
			logging.Panic(err)
			rightFileInfo, err := f.DirEntry.Info();
			logging.Panic(err)

			if leftFileInfo.ModTime().After(rightFileInfo.ModTime()) {
				f.Path = path
				f.DirEntry = file

				continue
			}
		}
	}
}
