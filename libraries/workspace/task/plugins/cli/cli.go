package cli

import (
	"context"
	"fmt"
	"github.com/satori/go.uuid"
	clil "github.com/walterjwhite/go-application/libraries/cli"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"os"
	"path/filepath"
)

func Execute(ctx context.Context, t *task.Task, captureScreenshots bool, cmd string, arguments ...string) {
	// TODO: support Env
	// Dir
	currentWorkingDirectory, err := os.Getwd()
	logging.Panic(err)

	id := uuid.Must(uuid.NewV4()).String()
	logDirectory := getLogDirectory(t, id)

	c := &clil.Command{
		Id:                 id,
		Dir:                currentWorkingDirectory,
		Cmd:                cmd,
		Args:               arguments,
		CaptureScreenshots: captureScreenshots,
		LogDirectory:       logDirectory,
	}

	c.Execute(ctx)

	t.WorkTreeConfig.Add(logDirectory)
	t.WorkTreeConfig.Commit(ctx, fmt.Sprintf("executed: %v", cmd))
	t.WorkTreeConfig.Push(ctx)
}

func getLogDirectory(t *task.Task, id string) string {
	return filepath.Join(t.Path, "logs", id)
}
