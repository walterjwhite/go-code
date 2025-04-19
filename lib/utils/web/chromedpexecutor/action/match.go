package action

import (
	"bytes"

	"github.com/andreyvit/locateimage"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"image/png"

	"context"
)

func Match(ctx context.Context, matchThreshold float64, referenceImageData []byte, x, y, width, height float64) []locateimage.Match {
	referenceImage, err := png.Decode(bytes.NewReader(referenceImageData))
	logging.Panic(err)

	screenshotData := TakeScreenshotOf(ctx, x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Panic(err)

	matches, err := locateimage.All(context.Background(), screenshotDataImg, referenceImage, matchThreshold)
	logging.Panic(err)

	log.Info().Msgf("matches: %v", matches)

	return matches
}
