package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/utils/workspace"
	"github.com/walterjwhite/go-application/libraries/utils/workspace/task"

	//"github.com/rs/zerolog/log"
	"flag"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/application"
	"github.com/walterjwhite/go-application/libraries/git/plugins/comment"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"os"
	"path/filepath"
)

func init() {
	application.Configure()
}

func main() {
	defer application.OnEnd()

	if len(flag.Args()) < 2 {
		logging.Panic(fmt.Errorf("Expecting action and names..."))
	}

	w := workspace.Get()
	taskNames := flag.Args()[1:]

	log.Info().Msgf("workspace: %v", w)
	log.Info().Msgf("workspace worktree: %v", w.WorkTreeConfig)

	switch flag.Args()[0] {
	case "new":
		for _, taskName := range taskNames {
			task.New(application.Context, w, taskName)
		}
	case "open":
		for _, taskName := range taskNames {
			t := getTask(w, taskName)
			t.Open(application.Context)
		}
	case "cancel":
		for _, taskName := range taskNames {
			t := getTask(w, taskName)
			t.Cancel(application.Context)
		}
	case "complete":
		for _, taskName := range taskNames {
			t := getTask(w, taskName)
			t.Complete(application.Context)
		}
	case "comment":
		if len(flag.Args()) != 3 {
			logging.Panic(fmt.Errorf("comment requires 2 arguments: taskName message"))
		}

		t := getTask(w, flag.Args()[1])
		logging.Panic(os.Chdir(t.GetPath()))

		comment.New(t.WorkTreeConfig, flag.Args()[2]).Write(application.Context)
	default:
		logging.Panic(fmt.Errorf("%v is not understood", flag.Args()[0]))
	}
}

func getTask(w *workspace.Workspace, taskName string) *task.Task {
	absFilePath := filepath.Join(w.WorkTreeConfig.Path, taskName)
	_, err := os.Stat(absFilePath)
	if os.IsNotExist(err) {
		logging.Panic(fmt.Errorf("Unable to initialize: %v (does *NOT* exist)", absFilePath))
	}

	return task.Initialize(application.Context, w, taskName)
}

/*
func do(callback func(ctx context.Context)) {
	for _, taskName := range flag.Args()[1:] {
		callback(ctx)
	}
}
*/
