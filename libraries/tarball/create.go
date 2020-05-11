package tarball

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io"
	"os"
	"path/filepath"
)

type tarball struct {
	w  io.Writer
	tw *tar.Writer
}

// TODO: support multiple compression schemes here (bzip2, gzip, lzma, etc.)
func Create(filename, path string, compress bool) {
	exists(path)
	mkdir(filename)

	f, err := os.Create(filename)
	logging.Panic(err)
	defer f.Close()

	t := &tarball{w: f}

	if compress {
		gw := gzip.NewWriter(f)
		defer gw.Close()

		t.w = gw
	}

	tw := tar.NewWriter(t.w)
	defer tw.Close()

	t.tw = tw

	logging.Panic(filepath.Walk(path, t.walk))
}

func mkdir(filename string) {
	parent := filepath.Dir(filename)
	_, err := os.Stat(parent)
	if os.IsNotExist(err) {
		logging.Panic(os.MkdirAll(parent, os.ModePerm))
	}
}

func exists(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logging.Panic(fmt.Errorf("%v does not exist, unable to create an archive of it.", path))
	}
}

// TODO: support excluding certain files or directories by name ...
func (t *tarball) walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	if info != nil && info.Mode()&os.ModeType == 0 {

		file, err := os.Open(path)
		logging.Panic(err)

		defer file.Close()
		stat, err := file.Stat()
		logging.Panic(err)

		header := new(tar.Header)
		header.Name = path
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()

		err = t.tw.WriteHeader(header)
		logging.Panic(err)

		_, err = io.Copy(t.tw, file)
		logging.Panic(err)
		return nil
	}

	log.Info().Msgf("skipping file: %v", path)
	return nil
}
