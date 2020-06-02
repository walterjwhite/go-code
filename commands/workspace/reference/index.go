package main

import (
	"github.com/walterjwhite/go-application/libraries/git/plugins/codesearch"
	"github.com/walterjwhite/go-application/libraries/workspace"

	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
)

func index() {
	if len(flag.Args()) == 2 {
		if *currentFlag {
			indexCurrent()
		} else {
			indexAll()
		}
	} else {
		indexSpecified()
	}
}

func indexSpecified() {
	for _, workspaceName := range flag.Args()[1:] {
		doIndex(workspace.Init(workspaceName))
	}
}

func indexAll() {
	wall := workspace.GetAll()
	for _, w := range wall {
		doIndex(w)
	}
}

func indexCurrent() {
	doIndex(workspace.Get())
}

func doIndex(w *workspace.Workspace) {
	codesearch.Index(application.Context, w.WorkTreeConfig)
}
