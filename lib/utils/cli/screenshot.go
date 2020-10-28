package cli

import (
	"github.com/walterjwhite/go/lib/utils/screenshot"

	"path/filepath"
)

func (c *Command) takeScreenshot(parentDirectory, name string) *screenshot.Instance {
	if c.CaptureScreenshots {
		filename := filepath.Join(parentDirectory, name)
		i := screenshot.Default(filename)

		i.Wait()
		return i
	}

	return nil
}
