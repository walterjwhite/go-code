package workspace

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/git"
	//"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Init(name string) *Workspace {
	loadProperties()

	w := &Workspace{Name: name}

	w.WorkTreeConfig = git.InitWorkTree(filepath.Join(Config.WorkspaceWorkPath, name))

	return w
}

func Get() *Workspace {
	loadProperties()

	currentWorkingDirectory, err := os.Getwd()
	logging.Panic(err)

	// remove workspace prefix
	prefixIndex := strings.Index(currentWorkingDirectory, Config.WorkspaceWorkPath)
	if prefixIndex == 0 {
		relativePath := strings.Replace(currentWorkingDirectory, Config.WorkspaceWorkPath, "", 1)
		parts := strings.Split(relativePath, string(os.PathSeparator))

		if len(parts) == 0 || len(relativePath) == 0 {
			logging.Panic(fmt.Errorf("GetWorkspace *MUST* be called within a workspace (in workspaces worktree, but not workspace)"))
		}

		//log.Info().Msgf("relative paths: %v (%v)", relativePath, parts)
		return &Workspace{Name: parts[1], WorkTreeConfig: git.InitWorkTree(filepath.Join(Config.WorkspaceWorkPath, parts[1]))}
	}

	logging.Panic(fmt.Errorf("GetWorkspace *MUST* be called within a workspace %d:%s", prefixIndex, currentWorkingDirectory))
	return nil
}

func GetAll() []*Workspace {
	loadProperties()

	files, err := ioutil.ReadDir(Config.WorkspaceWorkPath)
	logging.Panic(err)

	workspaces := make([]*Workspace, 0)
	for _, file := range files {
		if file.IsDir() {
			//logging.Panic(os.Chdir(filepath.Join(Config.WorkspaceWorkPath, file.Name())))
			// workspaces = append(workspaces, Get())
			workspaces = append(workspaces, Init(file.Name()))
		}
	}

	return workspaces
}
