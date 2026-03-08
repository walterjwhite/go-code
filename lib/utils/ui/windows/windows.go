package windows

import (
	"context"
	"os"
	"strconv"

	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"image"
)

const (
	defaultMatchThreshold = 0.80 // Default threshold for reliable image matching (0-1 scale)
	minMatchThreshold     = 0.50 // Minimum safe threshold to prevent false positives
	maxMatchThreshold     = 0.95 // Maximum threshold for strict matching
)

func getMatchThreshold() float64 {
	if envVal := os.Getenv("IMAGE_MATCH_THRESHOLD"); envVal != "" {
		if val, err := strconv.ParseFloat(envVal, 64); err == nil {
			if val >= minMatchThreshold && val <= maxMatchThreshold {
				return val
			}
		}
	}
	return defaultMatchThreshold
}

var matchThreshold = getMatchThreshold()

func (c *WindowsConf) IsLocked(ctx context.Context) (bool, error) {
	isPresent, err := c.IsStartMenuButtonPresent(ctx)
	return !isPresent, err
}

func (c *WindowsConf) IsStartMenuButtonPresent(ctx context.Context) (bool, error) {
	i, err := c.WindowsStartButtonMatcher(ctx)
	if err != nil {
		return false, err
	}

	return i.Matches()
}

func (c *WindowsConf) WindowsStartButtonMatcher(ctx context.Context) (*graphical.ImageMatch, error) {


	x, y, width, height, err := c.getScreenshotCoordinates(ctx)
	if err != nil {
		return nil, err
	}

	return &graphical.ImageMatch{Ctx: ctx, Image: c.getStartButtonImage(),
		MatchRegion: &graphical.MatchRegion{X: x, Y: y, Width: width, Height: height}, MatchThreshold: matchThreshold, Controller: c.Controller}, nil
}

func (c *WindowsConf) getStartButtonImage() image.Image {
	if c.Version == Windows10 {
		return windows10StartButtonImage
	}

	return windows11StartButtonImage
}

func (c *WindowsConf) getScreenshotCoordinates(ctx context.Context) (float64, float64, float64, float64, error) {
	size, err := action.GetWindowSize(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if c.Centered {
		return 0, float64(size.Height) - c.StartButtonHeight,
			float64(size.Width), c.StartButtonHeight, nil
	}

	return 0, float64(size.Height) - c.StartButtonHeight, c.StartButtonHeight * 2, c.StartButtonHeight, nil
}
