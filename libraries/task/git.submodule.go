package task

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (t *Task) initSubmodule(ctx context.Context, submodulePath string) {
	// initialize bare repository first
	//submoduleUri := t.doInitRemoteMirror(submodulePath)

	submoduleUri := t.prepareSubmoduleRepository(submodulePath)

	submoduleName := filepath.Join("active", submodulePath)

	cmd := runner.Prepare(ctx, "git", "submodule", "add", submoduleUri, submoduleName)
	cmd.Dir = gitSettings.WorkTreePath

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())

	_, err := t.w.Commit("initialized", &git.CommitOptions{Author: &object.Signature{Name: "Walter White", Email: "Walter.White@walterjwhite.com"}})
	logging.Panic(err)

	logging.Panic(t.git.Push(&git.PushOptions{}))
}

func (t *Task) prepareSubmoduleRepository(submodulePath string) string {
	dir, err := ioutil.TempDir("", "submodule")
	logging.Panic(err)

	defer os.RemoveAll(dir)

	tr, err := git.PlainInit(dir, false)
	logging.Panic(err)

	// commit something to the submodule
	_, err = os.Create(filepath.Join(dir, ".gitignore"))
	logging.Panic(err)

	tw, err := tr.Worktree()
	logging.Panic(err)

	_, err = tw.Add(".gitignore")
	logging.Panic(err)

	_, err = tw.Commit("initialized", &git.CommitOptions{Author: &object.Signature{Name: "Walter White", Email: "Walter.White@walterjwhite.com"}})
	logging.Panic(err)

	submoduleUri := getRemoteMirrorPath(submodulePath)
	_, err = git.PlainClone(submoduleUri, true, &git.CloneOptions{
		URL: dir,
	})
	logging.Panic(err)

	return submoduleUri
}
