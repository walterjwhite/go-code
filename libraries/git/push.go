package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Push(parentContext context.Context, projectDirectory string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	cmd := runner.Prepare(ctx, "git", "push")
	cmd.Dir = projectDirectory

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
