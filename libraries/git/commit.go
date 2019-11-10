package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Commit(parentContext context.Context, messageTemplate *string, commitMessage string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	_, err := runner.Run( /*application.Context*/ ctx, "git", "commit", "-m", FormatCommitMessage(messageTemplate, commitMessage))
	logging.Panic(err)
}
