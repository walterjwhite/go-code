package cli

import (
	"bytes"

	"github.com/satori/go.uuid"

	"context"

	"github.com/walterjwhite/go-application/libraries/application/logging"
	"io"
	"os"
	"time"

	"os/exec"

	"github.com/walterjwhite/go-application/libraries/io/yaml"
	"github.com/walterjwhite/go-application/libraries/time/timeout"
	"io/ioutil"
	"path/filepath"
	//"github.com/edwingeng/wuid/callback/wuid"
)

// creates the command object
// saves the command object (file, git, ES)
// executes the command
// captures output
// saves the updated command object (file, git, ES)

// TODO: compress the log file (lz4)
func (c *Command) Execute(ctx context.Context) {
	//g := NewWUID("default", nil)
	if len(c.Id) == 0 {
		c.Id = uuid.Must(uuid.NewV4()).String()
	}

	c.Date = time.Now()

	//_, filename := activity.Add(ctx, command, project)
	c.ctx = ctx

	c.doExecute()

	//activity.Update(ctx, command, project)
}

func (c *Command) doExecute() {
	if _, err := os.Stat(c.LogDirectory); os.IsNotExist(err) {
		logging.Panic(os.MkdirAll(c.LogDirectory, os.ModePerm))
	}

	sbefore := c.takeScreenshot(c.LogDirectory, "before")

	if c.TimeLimit != nil {
		timeout.Limit(c.run, c.TimeLimit, c.ctx)
	} else {
		c.run()
	}

	safter := c.takeScreenshot(c.LogDirectory, "after")

	// wait for images to be written out
	if sbefore != nil {
		sbefore.Wait()
		safter.Wait()
	}
}

func (c *Command) run() {
	cmd := exec.CommandContext(c.ctx, c.Cmd, c.Args...)
	cmd.Dir = c.Dir

	// TODO: convert from map to []string
	//cmd.Env = c.Env
	var outBuffer, errBuffer bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuffer)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuffer)

	logging.Panic(cmd.Run())

	// TODO: generalize this, don't want to write this twice
	//c.Stdout = outBuffer.String()
	//c.Stderr = errBuffer.String()

	//c.Status = cmd.ExitCode()
	c.CompletionDate = time.Now()

	// TODO: in some cases, we may not want to write to the FS
	// generalize this
	logging.Panic(ioutil.WriteFile(filepath.Join(c.LogDirectory, "stdout"), outBuffer.Bytes(), 0644))
	logging.Panic(ioutil.WriteFile(filepath.Join(c.LogDirectory, "stderr"), errBuffer.Bytes(), 0644))

	// write cmd to yaml
	yamlhelper.Write(c, filepath.Join(c.LogDirectory, "cli.yaml"))
}
