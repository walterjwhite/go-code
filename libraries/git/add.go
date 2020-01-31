package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Add(parentContext context.Context, projectDirectory string, filenames ...string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	filenames = append([]string{"add"}, filenames...)

	cmd := runner.Prepare(ctx, "git", filenames...)
	cmd.Dir = projectDirectory

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
