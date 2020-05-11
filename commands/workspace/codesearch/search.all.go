package main

import (
	"github.com/walterjwhite/go-application/libraries/workspace"
)

func doSearchAll() {
	for _, w := range workspace.GetAll() {
		searchWorkspace(w)
	}
}
