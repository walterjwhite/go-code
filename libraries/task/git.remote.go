package task

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

// TODO: this only supports local remotes
func (t *Task) initRemoteMirror() string {
	return t.doInitRemoteMirror("")
}

func (t *Task) doInitRemoteMirror(path string) string {
	remotePath := getRemoteMirrorPath(path) + ".git"

	_, err := os.Stat(remotePath)
	if os.IsNotExist(err) {
		_, err := git.PlainInit(remotePath, true)
		logging.Panic(err)
	}

	return remotePath
}

func getRemoteMirrorPath(path string) string {
	if len(path) == 0 {
		return gitSettings.RemoteUri
	}

	return filepath.Join(gitSettings.RemoteUri, path)
}
