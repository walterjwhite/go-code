package main

import (
	"errors"
	"flag"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/logging"
)

var (
	branchName = flag.String("GitRebaseBranchName", "", "Branch Name")
)

func init() {
	application.Configure()

	if len(*branchName) == 0 {
		logging.Panic(errors.New("Branch Name is required"))
	}
}

// TODO: integrate win10 / dbus notifications
func main() {
	r := &git.RebaseRequest{BranchName: *branchName}
	r.Rebase(application.Context)
}
