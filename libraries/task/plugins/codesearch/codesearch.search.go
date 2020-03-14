package codesearch

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
	"github.com/walterjwhite/go-application/libraries/task"
	"os"
)

// search indexes
func (c *Codesearch) Search(ctx context.Context, t *task.Task, pattern string) {
	c.doSearch(ctx, t, "", pattern)

	if len(c.GlobalIndexName) > 0 {
		c.doSearch(ctx, t, c.GlobalIndexName, pattern)
	}
}

func (c *Codesearch) doSearch(ctx context.Context, t *task.Task, indexName, pattern string) {
	indexPath, _ := c.getIndex(t, indexName)

	cmd := runner.Prepare(ctx, "csearch", pattern)
	cmd.Env = append(cmd.Env, fmt.Sprintf("CSEARCHINDEX=%v", indexPath))
	runner.WithWriter(cmd, os.Stdout)

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())
}
