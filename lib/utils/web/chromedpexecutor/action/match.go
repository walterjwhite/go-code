package action

import (
	"bytes"
	"github.com/andreyvit/locateimage"
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"image/png"

	"context"
	"errors"
	"os"
)

func Match(ctx context.Context, matchThreshold float64, referenceImageData []byte, x, y, width, height float64) *locateimage.Match {
	log.Debug().Msg("match start")

	referenceImage, err := png.Decode(bytes.NewReader(referenceImageData))
	logging.Panic(err)

	screenshotData := TakeScreenshotOf(ctx, x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Panic(err)

	if log.Debug().Enabled() {
		tempFile, err := ioutil.TempFile("/tmp", "match-*.png")
		logging.Panic(err)

		log.Debug().Msgf("capturing screenshot: %v", tempFile.Name())
		logging.Panic(os.WriteFile(tempFile.Name(), screenshotData, 0644))

		fullTempFile, err := ioutil.TempFile("/tmp", "full-*.png")
		logging.Panic(err)

		FullScreenshot(ctx, fullTempFile.Name())
	}

	match, err := locateimage.Find(context.Background(), locateimage.Convert(screenshotDataImg), locateimage.Convert(referenceImage), matchThreshold, locateimage.Fastest)
	/*
	   serr, ok := err.(*locateimage.ErrNotFound)

	   if err.(*locateimage.ErrNotFound); ok {
	     return nil
	   } else {
	     logging.Panic(err)
	   }
	*/
	target := locateimage.ErrNotFound
	if errors.As(err, &target) {
		log.Debug().Msg("match end - no matches")
		return nil
	} else {
		log.Debug().Msg("match end - err")
		logging.Panic(err)
	}

	log.Info().Msgf("matches: %v", match)
	return &match
}

func Matches(ctx context.Context, matchThreshold float64, referenceImageData []byte, x, y, width, height float64) []locateimage.Match {
	referenceImage, err := png.Decode(bytes.NewReader(referenceImageData))
	logging.Panic(err)

	screenshotData := TakeScreenshotOf(ctx, x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Panic(err)

	matches, err := locateimage.All(context.Background(), locateimage.Convert(screenshotDataImg), locateimage.Convert(referenceImage), matchThreshold)
	logging.Panic(err)

	log.Info().Msgf("matches: %v", matches)

	return matches
}
