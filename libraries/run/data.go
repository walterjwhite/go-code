package run

import (
	"os/exec"
)

type Application struct {
	Name        string
	Command     string
	Arguments   []string
	LogMatcher  string
	Environment []string

	command *exec.Cmd
}

type Instance struct {
	Applications []*Application
}
