package action

import (
	"bytes"
	"github.com/andreyvit/locateimage"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"image"
	"image/png"

	"context"
	"errors"
	"os"
)

func Match(ctx context.Context, matchThreshold float64, referenceImageData image.Image, x, y, width, height float64) *locateimage.Match {
	log.Debug().Msg("match start")

	screenshotData := TakeScreenshotOf(ctx, x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Warn(err, false, "Match")
	if err != nil {
		return nil
	}

	if log.Debug().Enabled() {
		tempFile, err := os.CreateTemp("", "region-*.png")
		logging.Warn(err, false, "Match.CreateTempFile-1")
		if err != nil {
			return nil
		}


		log.Debug().Msgf("capturing screenshot: %v", tempFile.Name())
		err = os.WriteFile(tempFile.Name(), screenshotData, 0644)
		logging.Warn(err, false, "Match.WriteFile")
		if err != nil {
			return nil
		}

		fullTempFile, err := os.CreateTemp("", "full-*.png")
		logging.Warn(err, false, "Match.CreateTempFile-2")
		if err != nil {
			return nil
		}

		FullScreenshot(ctx, fullTempFile.Name())
	}

	match, err := locateimage.Find(context.Background(), locateimage.Convert(screenshotDataImg), referenceImageData, matchThreshold, locateimage.Fastest)
	if errors.Is(err, locateimage.ErrNotFound) {
		log.Debug().Msg("match end - no matches")
		return nil
	} else {
		log.Debug().Msg("match end - err")
		logging.Warn(err, false, "Match.locateimage.Find")
		if err != nil {
			return nil
		}
	}

	log.Info().Msgf("matches: %v", match)
	return &match
}

func Matches(ctx context.Context, matchThreshold float64, referenceImageData image.Image, x, y, width, height float64) []locateimage.Match {
	screenshotData := TakeScreenshotOf(ctx, x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Panic(err)

	matches, err := locateimage.All(context.Background(), locateimage.Convert(screenshotDataImg), referenceImageData, matchThreshold)
	logging.Panic(err)

	log.Info().Msgf("matches: %v", matches)

	return matches
}
