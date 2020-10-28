package main

import (
	"flag"
	"github.com/walterjwhite/go/lib/application"
)

var (
	allFlag       = flag.Bool("all", false, "Search all tasks within all workspaces (default = false, search current task only)")
	workspaceFlag = flag.Bool("work", false, "Search all tasks within the current workspace (default = false, search current task only)")
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) == 0 {
		doIndex()
	} else {
		if *allFlag {
			doSearchAll()
		} else if *workspaceFlag {
			doSearchWorkspace()
		} else {
			doSearchTask()
		}
	}
}
