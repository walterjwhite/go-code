package main

import (
	"github.com/walterjwhite/go-application/libraries/workspace/task/plugins"
	"github.com/walterjwhite/go-application/libraries/workspace/task/plugins/codesearch"

	"github.com/walterjwhite/go-application/libraries/application"
)

func doIndex() {
	codesearch.Index(application.Context, plugins.Initialize(application.Context))
}
