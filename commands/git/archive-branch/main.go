package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	sourceBranchName = flag.String("GitArchiveBranchSourceBranchName", "", "Source Branch Name")
	branchName       = flag.String("GitArchiveBranchBranchName", "", "Branch Name")
)

func init() {
	application.Configure()

	if len(sourceBranchName) == 0 {
		logging.Panic(errors.New("Source Branch Name is required"))
	}

	if len(branchName) == 0 {
		logging.Panic(errors.New("Branch Name is required"))
	}
}

// TODO: integrate win10 / dbus notifications
func main() {
	r := &git.ArchiveRequest{BranchName: *branchName, SourceBranchName: *sourceBranchName}
	r.ArchiveBranch(application.Context)
}
