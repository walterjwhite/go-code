package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

func init() {
	application.Configure()
}

func main() {
	url, directory, submoduleUri, submodulePath := os.Args[1], os.Args[2], os.Args[3], os.Args[4]

	log.Info().Msgf("git clone %s %s", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url, Progress: os.Stdout,
	})
	logging.Panic(err)

	w, err := r.Worktree()
	logging.Panic(err)

	createSubmodule(application.Context, directory, submoduleUri, submodulePath)

	sms, err := w.Submodules()
	logging.Panic(err)

	log.Info().Msgf("submodules: %v", len(sms))
}

func createSubmodule(ctx context.Context, directory, submoduleUri, submodulePath string) {
	cmd := runner.Prepare(ctx, "git", "submodule", "add", submoduleUri, submodulePath)
	cmd.Dir = directory

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
