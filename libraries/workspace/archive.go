package workspace

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/git/plugins/archive"
	"github.com/walterjwhite/go-application/libraries/logging"

	"path/filepath"
)

func Archive(ctx context.Context, name string) {
	loadProperties()

	createArchive(name)
}

func createArchive(name string) {
	hExpand, err := homedir.Expand(Config.WorkspaceWorkPath)
	logging.Panic(err)

	wPath := filepath.Join(hExpand, name)

	// @see: https://gist.github.com/jonmorehouse/9060515
	fullPath, err := homedir.Expand(Config.ArchivePath)
	logging.Panic(err)

	// TODO: perhaps, embed a timestamp in the filename
	archivePath := filepath.Join(fullPath, name+".tar.gz")

	archive.Create(archivePath, wPath)
}
