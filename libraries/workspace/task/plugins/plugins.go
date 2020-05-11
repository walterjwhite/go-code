package plugins

import (
	"context"
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/property"
	"github.com/walterjwhite/go-application/libraries/typename"
	"github.com/walterjwhite/go-application/libraries/workspace"
	"github.com/walterjwhite/go-application/libraries/workspace/task"
	"os"
	"path/filepath"
)

// create task instance from current directory (look for .git)
func InitializeWithName(ctx context.Context) (*task.Task, string) {
	t := Initialize(ctx)

	name := ""

	if len(flag.Args()) >= 1 {
		name = flag.Args()[0]
	}

	return t, name
}

func Initialize(ctx context.Context) *task.Task {
	w := workspace.Get()

	currentWorkingDirectory, err := os.Getwd()
	logging.Panic(err)

	return InitializeTaskIn(ctx, w, currentWorkingDirectory)
}

func InitializeTaskIn(ctx context.Context, w *workspace.Workspace, path string) *task.Task {
	return task.Initialize(ctx, w, getTaskDirectory(path))
}

func getTaskDirectory(root string) string {
	gitDirectory := filepath.Join(root, ".git")

	_, err := os.Stat(gitDirectory)
	if os.IsNotExist(err) {
		parent := filepath.Dir(root)
		if len(parent) == 0 || parent == root {
			logging.Panic(fmt.Errorf("Unable to locate git directory"))
		}

		return getTaskDirectory(parent)
	}

	return root
}

func Configure(t *task.Task, name string, configuration interface{}) {
	propertyConfiguration := &property.Configuration{Path: GetConfiguration(t, name, configuration)}

	// TODO: support passing in prefix
	propertyConfiguration.Load(configuration, "")
}

func GetConfiguration(t *task.Task, name string, configuration interface{}) string {
	return filepath.Join(t.Path, "refs", name+"."+typename.Get(configuration))
}
