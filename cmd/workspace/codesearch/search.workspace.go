package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/utils/foreachfile"
	"github.com/walterjwhite/go/lib/utils/workspace"
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins"
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins/codesearch"

	"flag"
	"github.com/walterjwhite/go/lib/application"
	"os"
	"path/filepath"
)

type workspaceSearch struct {
	w *workspace.Workspace
}

func doSearchWorkspace() {
	searchWorkspace(workspace.Get())
}

func searchWorkspace(w *workspace.Workspace) {
	a := &workspaceSearch{w: w}

	foreachfile.Execute(a.w.GetWorkTreePath(), a.doSearch, ".gitignore")
}

func (a *workspaceSearch) doSearch(filePath string) {
	taskDirectory := filepath.Dir(filePath)

	// TODO: rather than panic in search, return an err, then recover here
	_, err := os.Stat(filepath.Join(taskDirectory, ".codesearch"))
	if err != nil && os.IsNotExist(err) {
		log.Warn().Msgf(".codesearch not initialized: %v", taskDirectory)
		return
	}

	t := plugins.InitializeTaskIn(application.Context, a.w, taskDirectory)
	codesearch.Search(application.Context, t, flag.Args()[0])
}
