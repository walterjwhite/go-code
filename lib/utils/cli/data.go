package cli

import (
	"context"
	"time"
)

type Command struct {
	Id string

	Date time.Time

	TimeLimit *time.Duration

	Env map[string]string `yaml:",omitempty"`
	Dir string            `yaml:",omitempty"`

	Cmd  string
	Args []string `yaml:",omitempty"`

	Tags []string `yaml:",omitempty"`

	// Command Output
	Stdout string `yaml:",omitempty"`
	Stderr string `yaml:",omitempty"`

	Status int

	// use a std filename
	//ScreenshotFilename string
	CaptureScreenshots bool

	CompletionDate time.Time

	LogDirectory string
	ctx          context.Context
}

type ScriptFile struct {
	Commands []Command

	DelayBetweenCommands *time.Duration
}

func (c *Command) DocumentId() string {
	return c.Id
}

func (c *Command) ActivityDateTime() time.Time {
	return c.Date
}

func (c *Command) Equals(record []string) bool {
	return c.Id == record[0]
}
