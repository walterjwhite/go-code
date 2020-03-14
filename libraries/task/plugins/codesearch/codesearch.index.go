package codesearch

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/task"
	"path/filepath"
)

// add/update
// TODO: support a global index too?
// TODO: add support for resetting the index
func (c *Codesearch) Index(ctx context.Context, t *task.Task) {
	c.doIndex(ctx, t, "")

	/*
		if len(c.GlobalIndexName) > 0 {
			doIndex(c.GlobalIndexName)
		}
	*/
}

func (c *Codesearch) doIndex(ctx context.Context, t *task.Task, indexName string) {
	indexPath, contentPath := c.getIndex(t, indexName)

	cmd := runner.Prepare(ctx, "cindex", contentPath)
	cmd.Env = append(cmd.Env, fmt.Sprintf("CSEARCHINDEX=%v", indexPath))

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}

func (c *Codesearch) getIndex(t *task.Task, indexName string) (string, string) {
	contentPath := t.Path
	indexPath := filepath.Join(t.Path, ".codesearch")
	// worktree + /.codesearch/tasks/<taskName> -> worktree/<status>/<taskName>
	// worktree + /.codesearch/common/<globalIndexName> -> worktree/<status>/<taskName>

	return indexPath, contentPath
}
