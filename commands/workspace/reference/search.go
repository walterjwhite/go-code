package main

import (
	"github.com/walterjwhite/go-application/libraries/git/plugins/codesearch"
	"github.com/walterjwhite/go-application/libraries/workspace"

	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
)

func search() {
	if len(flag.Args()) > 2 {
		for _, workspaceName := range flag.Args()[2:] {
			doSearch(workspace.Init(workspaceName))
		}
	} else {
		if *allFlag {
			wall := workspace.GetAll()
			for _, w := range wall {
				doSearch(w)
			}
		} else {
			doSearch(workspace.Get())
		}
	}
}

func doSearch(w *workspace.Workspace) {
	codesearch.Search(application.Context, w.WorkTreeConfig, flag.Args()[1])
}
