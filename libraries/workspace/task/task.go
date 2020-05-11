package task

import (
	"github.com/walterjwhite/go-application/libraries/workspace"
	"path/filepath"
)

/*
func init() {
	Git = &GitSettings{}
}
*/

func (t *Task) GetPath() string {
	return filepath.Join(workspace.Config.WorkspaceWorkPath, t.Workspace.Name, t.Path)
}
