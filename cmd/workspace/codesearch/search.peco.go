package main

import (
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins"
	"github.com/walterjwhite/go/lib/utils/workspace/task/plugins/codesearch"

	"flag"
	"github.com/walterjwhite/go/lib/application"

	"runtime"

	"context"

	"github.com/lestrrat-go/pdebug"
	"github.com/peco/peco"
	"github.com/peco/peco/internal/util"
)

func doSearchWithPeco() {
	//codesearch.Search(application.Context, plugins.Initialize(application.Context), flag.Args()[0])

	p := peco.New()
	// process stdin ...
	err := p.Run(ctx)
}
