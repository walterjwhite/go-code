package main

import (
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins"
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins/codesearch"

	"github.com/walterjwhite/go/lib/application"
)

func doIndex() {
	codesearch.Index(application.Context, plugins.Initialize(application.Context))
}
