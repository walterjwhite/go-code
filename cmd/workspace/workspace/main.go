package main

import (
	"github.com/walterjwhite/go/lib/utils/workspace"

	//"github.com/rs/zerolog/log"
	"context"
	"flag"
	"fmt"
	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) < 2 {
		logging.Panic(fmt.Errorf("Expecting action and workspace names..."))
	}

	workspaceNames := flag.Args()[1:]

	switch flag.Args()[0] {
	case "create":
		do(workspace.DoCreate, workspaceNames)
	case "archive":
		do(workspace.Archive, workspaceNames)
	default:
		logging.Panic(fmt.Errorf("%v is not understood", flag.Args()[0]))
	}
}

func do(callback func(ctx context.Context, name string), workspaceNames []string) {
	for _, workspaceName := range workspaceNames {
		callback(application.Context, workspaceName)
	}
}
