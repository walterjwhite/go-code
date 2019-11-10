package main

import (
	"errors"
	"flag"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	//	"github.com/walterjwhite/go-application/libraries/path"
	"github.com/walterjwhite/go-application/libraries/git"
	//	"github.com/walterjwhite/go-application/libraries/timestamp"
	//"os/exec"
)

// commit message format
// ticket-# message

var (
	commitMessageFormatFlag = flag.String("CommitMessageFormat", "%v %v", "Commit Message Format")
	commitMessage           string
)

func init() {
	application.Configure()

	commitMessage = flag.Args()[0]
	if len(commitMessage) == 0 {
		logging.Panic(errors.New("Commit Message is required"))
	} else {
		log.Info().Msgf("Using %v for the commit message", commitMessage)
	}
}

// TODO: integrate win10 / dbus notifications
func main() {
	//path.WithSessionDirectory("~/.audit/run/" + timestamp.Get())

	git.Commit(application.Context, commitMessageFormatFlag, commitMessage)
}
