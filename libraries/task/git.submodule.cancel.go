package task

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path/filepath"
)

func Cancel(ctx context.Context, submodulePath string) *Task {
	return changeStatus(ctx, submodulePath, "active", "canceled")
}

func (t *Task) move(ctx context.Context, submodulePath, oldStatus, newStatus string) {
	originalSubmoduleName := filepath.Join(oldStatus, submodulePath)
	newSubmoduleName := filepath.Join(newStatus, submodulePath)

	// this corrupts the repository ...
	/*
		_, err := t.w.Move(originalSubmoduleName, newSubmoduleName)
		logging.Panic(err)
	*/
	newTarget := filepath.Join(gitSettings.WorkTreePath, newSubmoduleName)
	parent := filepath.Dir(newTarget)
	_, err := os.Stat(parent)
	if os.IsNotExist(err) {
		logging.Panic(os.MkdirAll(parent, 0755))
	}

	cmd := runner.Prepare(ctx, "git", "mv", originalSubmoduleName, newSubmoduleName)
	cmd.Dir = gitSettings.WorkTreePath
	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())

	_, err = t.w.Commit(newStatus, &git.CommitOptions{Author: &object.Signature{Name: "Walter White", Email: "Walter.White@walterjwhite.com"}})
	logging.Panic(err)

	logging.Panic(t.git.Push(&git.PushOptions{}))
}
