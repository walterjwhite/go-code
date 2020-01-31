package git

import (
	"bufio"
	"bytes"
	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

type RebaseRequest struct {
	BranchName string

	// infer the source branch
	sourceBranchName string

	// return back to normal
	originalBranchName string

	//projectName string

	//patchFileName string

	successful bool
}

func (r *RebaseRequest) Rebase(ctx context.Context, projectDirectory string) {
	if !r.canOperate(ctx) {
		log.Error().Msg("Please commit your changes before attempting to rebase.")
	}

	r.sourceBranchName = GetSourceBranch(r.BranchName)
	r.originalBranchName = GetCurrentBranch(projectDirectory)

	// return back to normal
	defer r.restore(ctx)

	r.successful = false

	// ensure our source branch is up-to-date
	_, err := runner.Run(ctx, "git", "checkout", r.sourceBranchName)
	logging.Panic(err)

	_, err = runner.Run(ctx, "git", "pull")
	logging.Panic(err)

	_, err = runner.Run(ctx, "git", "checkout", r.BranchName)
	logging.Panic(err)

	_, err = runner.Run(ctx, "git", "rebase", r.sourceBranchName)
	logging.Panic(err)

	r.successful = true
}

func (r *RebaseRequest) canOperate(ctx context.Context) bool {
	var buffer bytes.Buffer

	cmd := runner.Prepare(ctx, "git", "status", "-s")
	runner.WithWriter(cmd, &buffer)
	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())

	scanner := bufio.NewScanner( /*output*/ &buffer)
	for scanner.Scan() {
		line := scanner.Text()

		if string(line[0]) == "A" || string(line[1]) == "M" || string(line[1]) == "D" || string(line[0:2]) == "??" {
			return false
		}
	}

	return true
}

func (r *RebaseRequest) restore(ctx context.Context) {
	if !r.successful {
		log.Warn().Msg("Reverting rebase as")
		_, err := runner.Run(ctx, "git", "rebase", "--abort")
		logging.Panic(err)
	}

	_, err := runner.Run(ctx, "git", "checkout", r.originalBranchName)
	logging.Panic(err)
}
