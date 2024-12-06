package cli

import (
	"bytes"

	"context"

	"io"
	"os"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"

	"os/exec"

	"path/filepath"

	"github.com/segmentio/ksuid"
	"github.com/walterjwhite/go-code/lib/io/yaml"
	"github.com/walterjwhite/go-code/lib/time/timeout"
)


func (c *Command) Execute(ctx context.Context) {
	if len(c.Id) == 0 {
		c.Id = ksuid.New().String()
	}

	c.Date = time.Now()

	c.ctx = ctx

	c.doExecute()

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

	if sbefore != nil {
		sbefore.Wait()
		safter.Wait()
	}
}

func (c *Command) run() {
	cmd := exec.CommandContext(c.ctx, c.Cmd, c.Args...)
	cmd.Dir = c.Dir

	var outBuffer, errBuffer bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuffer)
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuffer)

	logging.Panic(cmd.Run())


	c.CompletionDate = time.Now()

	logging.Panic(os.WriteFile(filepath.Join(c.LogDirectory, "stdout"), outBuffer.Bytes(), 0644))
	logging.Panic(os.WriteFile(filepath.Join(c.LogDirectory, "stderr"), errBuffer.Bytes(), 0644))

	yaml.Write(c, filepath.Join(c.LogDirectory, "cli.yaml"))
}
