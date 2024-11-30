package chromedpexecutor

import (
	"context"

	"math"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/chromedp/cdproto/emulation"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"os"
)

func Screenshot(s session.ChromeDPSession, filename string) {
	var buf []byte

	log.Info().Msgf("capturing screenshot: %v", filename)
	//logging.Panic(chromedp.Run(s.Context, fullScreenshot(90, &buf)))
	logging.Panic(chromedp.Run(s.Context(), chromedp.CaptureScreenshot(&buf)))

	log.Info().Msgf("took screenshot - writing to: %v", filename)

	logging.Panic(os.WriteFile(filename, buf, 0644))
	log.Info().Msgf("captured screenshot: %v", filename)
}

func FullScreenshot(quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			logging.Panic(err)

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			logging.Panic(err)

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			logging.Panic(err)

			return nil
		}),
	}
}
