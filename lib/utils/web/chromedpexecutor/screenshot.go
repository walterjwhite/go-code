package chromedpexecutor

import (
	"context"

	"math"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/chromedp/cdproto/emulation"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"os"
)

func Screenshot(ctx context.Context, filename string) {
	var buf []byte

	log.Debug().Msgf("capturing screenshot: %v", filename)
	logging.Panic(chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)))

	log.Debug().Msgf("took screenshot - writing to: %v", filename)

	logging.Panic(os.WriteFile(filename, buf, 0644))
	log.Debug().Msgf("captured screenshot: %v", filename)
}

func FullScreenshot(ctx context.Context, filename string) {
	var buf []byte
	logging.Panic(chromedp.Run(ctx, chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, _, _, _, err := page.GetLayoutMetrics().Do(ctx)
			logging.Panic(err)

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			log.Debug().Msgf("screen: [%d, %d]", width, height)
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			logging.Panic(err)

			log.Debug().Msgf("capture position: [%f, %f]", contentSize.X, contentSize.Y)
			log.Debug().Msgf("capture size: [%f, %f]", contentSize.Width, contentSize.Height)

			buf, err = page.CaptureScreenshot().
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
	}))

	logging.Panic(os.WriteFile(filename, buf, 0644))
	log.Debug().Msgf("captured screenshot: %v", filename)
}
