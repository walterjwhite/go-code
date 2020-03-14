package task

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func (t *Task) Commits() object.CommitIter {
	r, err := git.PlainOpen(t.Path)
	logging.Panic(err)

	ref, err := r.Head()
	logging.Panic(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	logging.Panic(err)

	return cIter
}
