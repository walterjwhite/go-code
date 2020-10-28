package run

import (
	"os/exec"
	"time"
)

type Application struct {
	Name string

	Command     string
	Arguments   []string
	LogMatcher  string
	Environment []string

	PortMonitorTimeout  time.Duration
	PortMonitorInterval time.Duration
	Port                int

	command *exec.Cmd
	session *Session

	//Files []string

	// each invocation will run out of a temporary directory
	//workDirectory
}

type Session struct {
	Path         string
	Applications []*Application
}
