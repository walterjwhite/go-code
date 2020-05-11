package main

import (
	"github.com/walterjwhite/go-application/libraries/workspace"

	"context"
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	allFlag = flag.Bool("a", false, "Search all references")
)

func init() {
	application.Configure()
	workspace.Name = "reference"
}

func main() {
	defer application.OnEnd()

	validateArgs()

	switch flag.Args()[0] {
	case "create":
		do(workspace.DoCreate, flag.Args()[1:])
	case "archive":
		do(workspace.Archive, flag.Args()[1:])
	// TODO: implement
	// case "commit":
	// 	do(workspace.Archive, flag.Args()[1:])
	case "index":
		index()
	case "search":
		search()
	default:
		logging.Panic(fmt.Errorf("%v is not understood", flag.Args()[0]))
	}
}

func do(callback func(ctx context.Context, name string), workspaceNames []string) {
	for _, workspaceName := range workspaceNames {
		callback(application.Context, workspaceName)
	}
}
