package task

import (
	"context"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/timeformatter/timestamp"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	commentPath        = "comments"
	commentPermissions = 0644
)

var (
	timestampConfiguration = &timestamp.Configuration{Template: "%d" + string(os.PathSeparator) + "%d" + string(os.PathSeparator) + "%d" + string(os.PathSeparator) + "%d.%d.%d"}
)

type comment struct {
	Message  string
	DateTime time.Time

	task *Task
}

func Comment(ctx context.Context, path, message string) *Task {
	t := initialize(ctx, path)

	c := &comment{Message: message, DateTime: time.Now()}
	commentPath := c.write()

	_, err := t.w.Add(commentPath)
	logging.Panic(err)

	_, err = t.w.Commit(message, &git.CommitOptions{})
	logging.Panic(err)

	logging.Panic(t.git.Push(&git.PushOptions{}))

	return t
}

func (c *comment) write() string {
	commentPath := c.absolute()
	logging.Panic(ioutil.WriteFile(commentPath, []byte(c.Message), commentPermissions))

	return commentPath
}

func (c *comment) absolute() string {
	return filepath.Join(c.task.Path, commentPath, c.relative())
}

func (c *comment) relative() string {
	return timestampConfiguration.Format(c.DateTime)
}
