package git

import (
	"context"
	"time"

	"github.com/tcnksm/go-gitconfig"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type commit struct {
	Message        string
	WorkTreeConfig *WorkTreeConfig
}

// TODO: reintroduce timeouts ...
func (c *WorkTreeConfig) Commit(parentCtx context.Context, commitMessage string) {
	commitConfig := &commit{Message: commitMessage, WorkTreeConfig: c}

	d := 30 * time.Second

	timeout.Limit(commitConfig.doCommit, &d, parentCtx)
}

func (c *commit) doCommit() {
	// TODO: get the signature from the environment
	username, err := gitconfig.Username()
	logging.Panic(err)

	email, err := gitconfig.Email()
	logging.Panic(err)

	_, err = c.WorkTreeConfig.W.Commit(c.Message, &git.CommitOptions{Author: &object.Signature{Name: username, Email: email, When: time.Now()}})

	logging.Panic(err)
}
