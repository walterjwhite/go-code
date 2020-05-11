package git

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
)

func InitCurrentWDWorkTree() *WorkTreeConfig {
	wd, err := os.Getwd()
	logging.Panic(err)

	return InitWorkTree(wd)
}

func InitWorkTreeIn(path string) *WorkTreeConfig {
	return InitWorkTree(getGitDir(path))
}

func getGitDir(path string) string {
	// check if .git is present
	gitDir := filepath.Join(path, ".git")
	_, err := os.Stat(gitDir)
	if os.IsNotExist(err) {
		parent := filepath.Dir(path)
		if parent == path {
			logging.Panic(fmt.Errorf("Unable to locate git directory in: %v", path))
		}

		return getGitDir(parent)
	}

	return path
}

func InitWorkTree(path string) *WorkTreeConfig {
	return doInit(path, filepath.Join(path, ".git"), false)
}

func InitBare(path string) *WorkTreeConfig {
	return doInit(path, path, true)
}

func doInit(path, gitDir string, bare bool) *WorkTreeConfig {
	expandedPath, err := homedir.Expand(path)
	logging.Panic(err)

	c := &WorkTreeConfig{Path: expandedPath}

	_, err = os.Stat(gitDir)
	if os.IsNotExist(err) {
		c.doInitRepository(git.PlainInit(expandedPath, bare))
	} else {
		c.doInitRepository(git.PlainOpen(expandedPath))
	}

	if !bare {
		c.doInitWorkTree()
	}

	return c
}

func (c *WorkTreeConfig) doInitRepository(r *git.Repository, err error) {
	logging.Panic(err)
	c.R = r
}

func (c *WorkTreeConfig) doInitWorkTree() {
	w, err := c.R.Worktree()
	logging.Panic(err)

	c.W = w
}
