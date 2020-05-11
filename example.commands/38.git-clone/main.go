package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"gopkg.in/src-d/go-git.v4"
)

var (
	urlFlag       = flag.String("Url", "ssh://git@localhost:/projects/active/go/walterjwhite.git", "Git URL")
	directoryFlag = flag.String("Directory", "/tmp/git-clone", "Workspace Directory")
)

func init() {
	application.Configure()
}

func main() {

	r, err := git.PlainClone(*directoryFlag, false, &git.CloneOptions{URL: *urlFlag})
	logging.Panic(err)

	ref, err := r.Head()
	logging.Panic(err)

	log.Info().Msgf("HEAD hash: %v", ref.Hash())
}
