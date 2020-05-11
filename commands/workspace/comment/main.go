package main

import (
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/git/plugins/comment"
	"github.com/walterjwhite/go-application/libraries/logging"
	"os"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) != 1 {
		logging.Panic(fmt.Errorf("Expecting comment only"))
	}

	wd, err := os.Getwd()
	logging.Panic(err)

	w := git.InitWorkTreeIn(wd)
	comment.New(w, flag.Args()[0]).Write(application.Context)
}
