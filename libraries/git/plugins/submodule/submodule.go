package submodule

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/git"
	//"github.com/walterjwhite/go-application/libraries/git/plugins/remote"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Add(ctx context.Context, w *git.WorkTreeConfig, submoduleUri, submoduleName string) {
	log.Debug().Msgf("submoduleUri / name: %v / %v", submoduleUri, submoduleName)
	cmd := runner.Prepare(ctx, "git", "submodule", "add", submoduleUri, submoduleName)
	cmd.Dir = w.Path

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}

// this corrupts the repository ...
/*
	_, err := t.w.Move(originalSubmoduleName, newSubmoduleName)
	logging.Panic(err)
*/
func Move(ctx context.Context, w *git.WorkTreeConfig, source, target string) {
	log.Debug().Msgf("dir: %v", w.Path)

	cmd := runner.Prepare(ctx, "git", "mv", source, target)
	cmd.Dir = w.Path

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}

func AtomicMove(ctx context.Context, w *git.WorkTreeConfig, source, target, commitMessage string) {
	_, err := os.Stat(filepath.Join(w.Path, source))
	if os.IsNotExist(err) {
		logging.Panic(fmt.Errorf("Unable to move: %v -> %v as source (%v) does not exist", source, target, source))
	}
	prepareTarget(w, target)

	Move(ctx, w, source, target)

	w.Commit(ctx, commitMessage)
	w.Push(ctx)
}

func prepareTarget(w *git.WorkTreeConfig, target string) {
	absoluteTarget := filepath.Join(w.Path, target)
	parent := filepath.Dir(absoluteTarget)
	_, err := os.Stat(parent)
	if os.IsNotExist(err) {
		logging.Panic(os.MkdirAll(parent, 0755))
	}
}

func AtomicAdd(ctx context.Context, w *git.WorkTreeConfig, submodulePath, submoduleName, mirrorPath string) {
	submoduleUri := prepareSubmoduleRepository(ctx, submodulePath, mirrorPath)
	Add(ctx, w, submoduleUri, submoduleName)

	w.Commit(ctx, fmt.Sprintf("initialized - %v", submoduleName))
	w.Push(ctx)
}

func prepareSubmoduleRepository(ctx context.Context, submodulePath, mirrorPath string) string {
	dir, err := ioutil.TempDir("", "submodule")
	logging.Panic(err)

	defer os.RemoveAll(dir)

	w := git.InitWorkTree(dir)

	_, err = os.Create(filepath.Join(dir, ".gitignore"))
	logging.Panic(err)

	w.Add(".gitignore")
	w.Commit(ctx, "submodule - initialized")

	//return remote.Init(w, mirrorPath, submodulePath)
	mw := &git.WorkTreeConfig{Path: filepath.Join(mirrorPath, submodulePath)}

	mw.Mirror(ctx, dir)

	return mw.Path
}
