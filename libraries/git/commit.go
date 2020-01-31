package git

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Commit(parentContext context.Context, projectDirectory string, messageTemplate *string, commitMessage string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	log.Info().Msgf("args: %v / %v", *messageTemplate, commitMessage)
	cmd := runner.Prepare(ctx, "git", "commit", "-m", FormatCommitMessage(projectDirectory, messageTemplate, commitMessage))
	cmd.Dir = projectDirectory

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
