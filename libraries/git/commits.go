package git

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func (c *WorkTreeConfig) Commits() object.CommitIter {
	ref, err := c.R.Head()
	logging.Panic(err)

	cIter, err := c.R.Log(&git.LogOptions{From: ref.Hash()})
	logging.Panic(err)

	return cIter
}
