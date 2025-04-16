package citrix

import (
	"bytes"
	"github.com/andreyvit/locateimage"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"image/png"

	"context"
)

func (i *Instance) takeScreenshot(x, y, width, height float64) []byte {
	log.Info().Msgf("taking screenshot: [%f, %f] [%f, %f]", x, y, width, height)

	var buf []byte
	logging.Panic(chromedp.Run(i.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		buf, err = page.CaptureScreenshot().
			WithClip(&page.Viewport{

				X:      x,
				Y:      y,
				Width:  width,
				Height: height,

				Scale: 1,
			}).Do(ctx)
		return err
	})))
	return buf
}

func (i *Instance) match(matchThreshold float64, referenceImageData []byte, x, y, width, height float64) []locateimage.Match {
	referenceImage, err := png.Decode(bytes.NewReader(referenceImageData))
	logging.Panic(err)

	screenshotData := i.takeScreenshot(x, y, width, height)
	screenshotDataImg, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Panic(err)

	matches, err := locateimage.All(context.Background(), screenshotDataImg, referenceImage, lockScreenMatchThreshold)
	logging.Panic(err)

	log.Info().Msgf("matches: %v", matches)

	return matches
}
