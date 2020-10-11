package run

import (
	"os/exec"
	"time"
)

type Application struct {
	Name        string
	Command     string
	Arguments   []string
	LogMatcher  string
	Environment []string

	PortMonitorTimeout  time.Duration
	PortMonitorInterval time.Duration
	Port                int
	command             *exec.Cmd

	//Files []string

	// each invocation will run out of a temporary directory
	//workDirectory
}

type Instance struct {
	Applications []*Application
}
