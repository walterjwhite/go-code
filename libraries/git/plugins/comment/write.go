package comment

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"

	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Comment) Write(ctx context.Context) {
	commentPath := c.doWrite()

	c.WorkTreeConfig.Add(commentPath)
	c.WorkTreeConfig.Commit(ctx, c.Message)

	c.WorkTreeConfig.Push(ctx)
}

func (c *Comment) doWrite() string {
	commentPath := c.absolute()

	logging.Panic(os.MkdirAll(filepath.Dir(commentPath), os.ModePerm))
	logging.Panic(ioutil.WriteFile(commentPath, []byte(c.Message), commentPermissions))

	relativePath, err := filepath.Rel(c.WorkTreeConfig.Path, commentPath)
	logging.Panic(err)

	return relativePath
}

func (c *Comment) absolute() string {
	log.Info().Msgf("path: %v", c.WorkTreeConfig.Path)
	return filepath.Join(c.WorkTreeConfig.Path, commentPath, c.relative())
}

func (c *Comment) relative() string {
	return timestampConfiguration.Format(c.DateTime)
}
