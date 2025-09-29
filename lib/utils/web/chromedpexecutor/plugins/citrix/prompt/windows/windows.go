package windows

import (
	"context"

	"github.com/andreyvit/locateimage"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action/windows"

	"image"
)

const (
	matchThreshold = 0.04
)

func (c *WindowsConf) IsLocked(ctx context.Context) (bool, error) {
	unlocked, err := c.IsWindowsStartButtonVisible(ctx)
	return !unlocked, err
}

func (c *WindowsConf) IsWindowsStartButtonVisible(ctx context.Context) (bool, error) {
	if c.AutomaticallyHides {
		err := windows.ToggleWindowsStartMenu(ctx)
		logging.Warn(err, false, "ToggleWindowsStartMenu")
		if err != nil {
			return false, err
		}

		defer func() {
			err := windows.ToggleWindowsStartMenu(ctx)
			logging.Warn(err, false, "ToggleWindowsStartMenu - defer - IsWindowsStartButtonVisible")
		}()
	}

	x, y, width, height, err := c.getScreenshotCoordinates(ctx)
	if err != nil {
		return false, err
	}

	match := action.Match(ctx, matchThreshold, c.getStartButtonImage(), x, y, width, height)

	if match != nil {
		log.Info().Msgf("found at %v, similarity = %.*f%%", match.Rect, locateimage.SimilarityDigits-2, 100*match.Similarity)
	} else {
		log.Warn().Msg("no match found")
	}

	return match != nil, nil
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
