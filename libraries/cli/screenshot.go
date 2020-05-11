package cli

import (
	"github.com/walterjwhite/go-application/libraries/screenshot"

	"path/filepath"
)

func (c *Command) takeScreenshot(parentDirectory, name string) *screenshot.Instance {
	if c.CaptureScreenshots {
		filename := filepath.Join(parentDirectory, name)
		return screenshot.Take(filename)
	}

	return nil
}
