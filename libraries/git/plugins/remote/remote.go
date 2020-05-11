package remote

import (
	"github.com/walterjwhite/go-application/libraries/logging"

	"github.com/walterjwhite/go-application/libraries/git"
	"gopkg.in/src-d/go-git.v4/config"

	"os"
	"path/filepath"
)

var (
	RemoteName = "origin"
)

// TODO: this only supports local remotes
func Init(w *git.WorkTreeConfig, mirrorPath, path string) string {
	remotePath := getRemoteMirrorPath(mirrorPath, path) + ".git"

	_, err := os.Stat(remotePath)
	if os.IsNotExist(err) {
		git.InitBare(remotePath)
	}

	Add(w, remotePath)

	return remotePath
}

func getRemoteMirrorPath(mirrorPath, path string) string {
	return filepath.Join(mirrorPath, path)
}

func Add(w *git.WorkTreeConfig, remotePath string) {
	if isSetupRemote(w) {
		setupRemote(w, remotePath)
	}
}

func isSetupRemote(w *git.WorkTreeConfig) bool {
	remotes, err := w.R.Remotes()
	logging.Panic(err)

	return len(remotes) == 0
}

func setupRemote(w *git.WorkTreeConfig, remotePath string) {
	remoteConfig := &config.RemoteConfig{Name: RemoteName, URLs: []string{remotePath}}
	_, err := w.R.CreateRemote(remoteConfig)
	logging.Panic(err)
}
