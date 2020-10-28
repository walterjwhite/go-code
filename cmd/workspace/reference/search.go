package main

import (
	"github.com/walterjwhite/go/lib/git/plugins/codesearch"
	"github.com/walterjwhite/go/lib/utils/workspace"

	"flag"
	"github.com/walterjwhite/go/lib/application"
)

func search() {
	if len(flag.Args()) == 2 {
		if *currentFlag {
			searchCurrent()
		} else {
			searchAll()
		}
	} else {
		searchSpecified()
	}
}

func searchSpecified() {
	for _, workspaceName := range flag.Args()[2:] {
		doSearch(workspace.Init(workspaceName))
	}
}

func searchAll() {
	wall := workspace.GetAll()
	for _, w := range wall {
		doSearch(w)
	}
}

func searchCurrent() {
	doSearch(workspace.Get())
}

func doSearch(w *workspace.Workspace) {
	codesearch.Search(application.Context, w.WorkTreeConfig, flag.Args()[1])
}
