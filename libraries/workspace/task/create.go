package task

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/git/plugins/submodule"
	//"github.com/walterjwhite/go-application/libraries/property"
	"github.com/walterjwhite/go-application/libraries/workspace"
	"path/filepath"
)

// create a new task at the specificed path
func New(ctx context.Context, w *workspace.Workspace, path string) *Task {
	loadProperties()
	t := &Task{Path: filepath.Join(Open, path), Workspace: w}

	submodule.AtomicAdd(ctx, w.WorkTreeConfig, path, t.Path, workspace.Config.WorkspaceRemotePath)

	t.WorkTreeConfig = git.InitWorkTree(t.GetPath())
	return t
}

func Initialize(ctx context.Context, w *workspace.Workspace, path string) *Task {
	loadProperties()
	t := &Task{Path: path, Workspace: w}

	log.Info().Msgf("opening: %v", t.Path)
	t.WorkTreeConfig = git.InitWorkTree(t.GetPath())

	return t
}

func loadProperties() {
	//property.Load(Git, "")
}
