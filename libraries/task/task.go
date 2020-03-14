package task

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/property"
	"gopkg.in/src-d/go-git.v4"
	//"path/filepath"
)

type Status string

/*
const (
	active    = "active"
	cancelled = "cancelled"
	pending   = "pending"
	completed = "completed"
)
*/

// each task gets its own repository
type Task struct {
	Path string

	// redundant, this would just be the submodules
	//Children []Task

	// redundant, we can just query the repository history
	//Commits []Commit
	Comments []*comment

	git *git.Repository
	w   *git.Worktree
}

func init() {
	gitSettings = &GitSettings{}
	property.Load(gitSettings, "")
}

// create a new task at the specificed path
func New(ctx context.Context, path string) *Task {
	t := &Task{Path: gitSettings.WorkTreePath}

	t.initRemoteMirror()
	t.initWorktree()

	t.initSubmodule(ctx, path)

	return t
}

// updates the task status
func changeStatus(ctx context.Context, submodulePath, oldStatus, status string) *Task {
	t := &Task{Path: gitSettings.WorkTreePath}
	t.initWorktree()

	//t := initialize(ctx, filepath.Join(gitSettings.WorkTreePath, oldStatus, submodulePath))
	t.move(ctx, submodulePath, oldStatus, status)

	return t
}

func initialize(ctx context.Context, path string) *Task {
	t := &Task{Path: path}

	log.Info().Msgf("opening: %v", t.Path)
	r, err := git.PlainOpen(t.Path)
	logging.Panic(err)

	t.git = r

	w, err := r.Worktree()
	logging.Panic(err)

	t.w = w

	return t
}
