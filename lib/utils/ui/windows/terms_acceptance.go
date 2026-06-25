package windows

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"image"
)

const (
	w10TermsAcceptanceButtonWidth = 300
	w11TermsAcceptanceButtonWidth = 75
)

func (c *WindowsConf) IsTermsAcceptanceButtonVisible(ctx context.Context) (bool, error) {
	x, y, width, height, err := c.getTermsAcceptanceButtonCoordinates(ctx)
	if err != nil {
		return false, err
	}

	i := graphical.ImageMatch{Ctx: ctx, Image: c.getTermsAcceptanceButtonImage(),
		MatchRegion: &graphical.MatchRegion{X: x, Y: y, Width: width, Height: height}, MatchThreshold: matchThreshold, Controller: c.Controller}

	return i.Matches()
}

func (c *WindowsConf) getTermsAcceptanceButtonImage() image.Image {
	if c.Version == Windows10 {
		return windows10TermsAcceptanceButtonImage
	}

	return windows11TermsAcceptanceButtonImage
}

func (c *WindowsConf) getTermsAcceptanceButtonCoordinates(ctx context.Context) (float64, float64, float64, float64, error) {
	size, err := action.GetWindowSize(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if !c.Centered {
		log.Debug().Msg("left aligned")
		return float64(size.Width)/2.0 - w10TermsAcceptanceButtonWidth, float64(size.Height) / 2.0, float64(w10TermsAcceptanceButtonWidth), float64(size.Height) / 4.0, nil
	}

	log.Debug().Msg("center aligned")
	return float64(size.Width)/2.0 - w11TermsAcceptanceButtonWidth, float64(size.Height) / 2.0, float64(w11TermsAcceptanceButtonWidth) * 2.0, float64(size.Height) / 4.0, nil
}
