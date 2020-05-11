package task

import (
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/workspace"
)

type Status string

const (
	Open      = "open"
	Cancelled = "cancelled"
	Completed = "completed"
)

// each task gets its own repository
type Task struct {
	Workspace *workspace.Workspace

	Path string

	WorkTreeConfig *git.WorkTreeConfig
}
