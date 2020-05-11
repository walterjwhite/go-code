package workspace

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-application/libraries/git"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/property"
)

type Workspace struct {
	Name string

	WorkTreeConfig *git.WorkTreeConfig

	// if this is publicly accessible on github (or another code hosting site)
	RemoteUri string
}

type Configuration struct {
	// ~/workspaces
	WorkspaceWorkPath   string
	WorkspaceRemotePath string

	// ~/.workspaces/archives
	ArchivePath string
}

var (
	Config *Configuration
	Name   string
	Prefix string
)

func init() {
	Name = "workspace"
	Prefix = ""
}

func loadProperties() {
	if Config == nil {
		Config = getDefault()
		property.Load(Config, Prefix)

		// expand paths ...
		expanded, err := homedir.Expand(Config.WorkspaceWorkPath)
		logging.Panic(err)
		Config.WorkspaceWorkPath = expanded

		expanded, err = homedir.Expand(Config.WorkspaceRemotePath)
		logging.Panic(err)
		Config.WorkspaceRemotePath = expanded

		expanded, err = homedir.Expand(Config.ArchivePath)
		logging.Panic(err)
		Config.ArchivePath = expanded
	}
}

func getDefault() *Configuration {
	return &Configuration{
		WorkspaceWorkPath:   fmt.Sprintf("~/%ss", Name),
		WorkspaceRemotePath: fmt.Sprintf("/tmp/%ss.mirror", Name),
		ArchivePath:         fmt.Sprintf("/tmp/%ss.archives", Name)}
}
