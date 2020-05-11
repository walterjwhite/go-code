package workspace

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/git/plugins/remote"
	"path/filepath"
	"strings"
)

func Create(ctx context.Context, name string) *Workspace {
	loadProperties()

	var w *Workspace
	if isRemote(name) {
		w = doRemote(ctx, name)
	} else {
		w = doLocal(name)
	}

	remote.Init(w.WorkTreeConfig, Config.WorkspaceRemotePath, name)

	return w
}

func isRemote(name string) bool {
	return strings.Contains(name, ",")
}

func doRemote(ctx context.Context, name string) *Workspace {
	parts := strings.Split(name, ",")
	w := &Workspace{Name: parts[0], RemoteUri: parts[1]}
	w.WorkTreeConfig = &git.WorkTreeConfig{Path: w.GetWorkTreePath()}
	w.WorkTreeConfig.Clone(ctx, w.RemoteUri)

	return w
}

func doLocal(name string) *Workspace {
	w := &Workspace{Name: name}
	w.WorkTreeConfig = git.InitWorkTree(w.GetWorkTreePath())

	return w
}

func (w *Workspace) GetWorkTreePath() string {
	return filepath.Join(Config.WorkspaceWorkPath, w.Name)
}

func DoCreate(ctx context.Context, name string) {
	Create(ctx, name)
}
