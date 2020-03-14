package task

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"os"
)

func (t *Task) initWorktree() {
	_, err := os.Stat(gitSettings.WorkTreePath)
	if os.IsNotExist(err) {
		r, err := git.PlainInit(gitSettings.WorkTreePath, false)
		logging.Panic(err)

		t.git = r

		w, err := r.Worktree()
		logging.Panic(err)

		t.w = w

		// setup remote
		remoteConfig := &config.RemoteConfig{Name: "origin", URLs: []string{gitSettings.RemoteUri}}
		_, err = r.CreateRemote(remoteConfig)
		logging.Panic(err)
	} else {
		r, err := git.PlainOpen(gitSettings.WorkTreePath)
		logging.Panic(err)

		t.git = r

		w, err := r.Worktree()
		logging.Panic(err)

		t.w = w
	}
}
