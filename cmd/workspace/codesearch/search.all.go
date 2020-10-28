package main

import (
	"github.com/walterjwhite/go/lib/utils/workspace"
)

func doSearchAll() {
	for _, w := range workspace.GetAll() {
		searchWorkspace(w)
	}
}
