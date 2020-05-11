package archive

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/tarball"
	"os"
)

func Create(archivePath, path string) {
	// w := git.InitWorkTree(path)

	tarball.Create(archivePath, path, true)
	delete(path)
}

func delete(path string) {
	logging.Panic(os.RemoveAll(path))
}
