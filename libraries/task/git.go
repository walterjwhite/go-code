package task

type GitSettings struct {
	// location of mirror, ie. https://github.com/walterjwhite/tasks.git
	// submodules would be,    https://github.com/walterjwhite/tasks/<status>/path.git
	RemoteUri string

	// local workspace, ie. ~/tasks
	WorkTreePath string
}

var (
	gitSettings *GitSettings
)
