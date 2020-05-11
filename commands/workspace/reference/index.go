package main

import (
	"github.com/walterjwhite/go-application/libraries/git/plugins/codesearch"
	"github.com/walterjwhite/go-application/libraries/workspace"

	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
)

func index() {
	if len(flag.Args()) > 1 {
		for _, workspaceName := range flag.Args()[1:] {
			doIndex(workspace.Init(workspaceName))
		}
	} else {
		if *allFlag {
			wall := workspace.GetAll()
			for _, w := range wall {
				doIndex(w)
			}
		} else {
			doIndex(workspace.Get())
		}
	}
}

func doIndex(w *workspace.Workspace) {
	codesearch.Index(application.Context, w.WorkTreeConfig)
}
