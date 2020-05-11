package git

import (
	"gopkg.in/src-d/go-git.v4"
)

type WorkTreeConfig struct {
	Path string

	R *git.Repository
	W *git.Worktree
}
