package windows

import (
	"context"

	"github.com/andreyvit/locateimage"

	"github.com/rs/zerolog/log"

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

	match := action.Match(ctx, matchThreshold, c.getTermsAcceptanceButtonImage(), x, y, width, height)

	if match != nil {
		log.Info().Msgf("found at %v, similarity = %.*f%%", match.Rect, locateimage.SimilarityDigits-2, 100*match.Similarity)
	} else {
		log.Warn().Msg("no match found")
	}

	return match != nil, nil
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

	if c.Version == Windows10 {
		log.Debug().Msg("windows 10 - left alignted")
		return float64(size.Width)/2.0 - w10TermsAcceptanceButtonWidth, float64(size.Height) / 2.0, float64(w10TermsAcceptanceButtonWidth), float64(size.Height) / 4.0, nil
	}

	log.Debug().Msg("windows 11 - center alignted")
	return float64(size.Width)/2.0 - w11TermsAcceptanceButtonWidth, float64(size.Height) / 2.0, float64(w11TermsAcceptanceButtonWidth) * 2.0, float64(size.Height) / 4.0, nil
}
