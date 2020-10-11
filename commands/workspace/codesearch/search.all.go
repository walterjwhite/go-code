package main

import (
	"github.com/walterjwhite/go-application/libraries/utils/workspace"
)

func doSearchAll() {
	for _, w := range workspace.GetAll() {
		searchWorkspace(w)
	}
}
