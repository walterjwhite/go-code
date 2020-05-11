package git

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeout"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

type cloneConfig struct {
	RemoteUri string
	Directory string

	Repository *git.Repository
	Bare       bool
}

func (c *WorkTreeConfig) Mirror(parentCtx context.Context, remoteUri string) {
	c.doClone(parentCtx, &cloneConfig{RemoteUri: remoteUri, Directory: c.Path, Bare: true})
}

func (c *WorkTreeConfig) Clone(parentCtx context.Context, remoteUri string) {
	c.doClone(parentCtx, &cloneConfig{RemoteUri: remoteUri, Directory: c.Path})
}

func (c *WorkTreeConfig) doClone(parentCtx context.Context, clone *cloneConfig) {
	log.Debug().Msgf("git clone %s %s", clone.RemoteUri, clone.Directory)

	d := 30 * time.Second
	timeout.Limit(clone.doTimeConstrainedClone, &d, parentCtx)

	c.R = clone.Repository

	if !clone.Bare {
		c.doInitWorkTree()
	}
}

func (c *cloneConfig) doTimeConstrainedClone() {
	r, err := git.PlainClone(c.Directory, c.Bare, getCloneOptions(c.RemoteUri))
	logging.Panic(err)

	c.Repository = r
}

func getCloneOptions(remoteUri string) *git.CloneOptions {
	o := &git.CloneOptions{URL: remoteUri}

	if log.Debug().Enabled() {
		o.Progress = os.Stdout
	}

	return o
}
