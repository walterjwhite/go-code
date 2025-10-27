package action

import (
	"context"

	"fmt"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"os"
)

func ScreenshotIfDebug(ctx context.Context, name_template string, args ...interface{}) {
	ScreenshotIf(ctx, log.Debug(), name_template, args...)
}

func ScreenshotIf(ctx context.Context, logLevel *zerolog.Event, name_template string, args ...interface{}) {
	if logLevel.Enabled() {
		Screenshot(ctx, fmt.Sprintf(name_template, args...))
	}
}

func Screenshot(pctx context.Context, filename string) {
	log.Debug().Msgf("capturing screenshot: %v", filename)

	log.Debug().Msgf("took screenshot - writing to: %v", filename)

	buf, err := TakeScreenshot(pctx)
	if err != nil {
		log.Warn().Msgf("error taking screenshot: %v", err)
		return
	}

	logging.Warn(os.WriteFile(filename, buf, 0644), false, "Screenshot.WriteFile")
	log.Debug().Msgf("captured screenshot: %v", filename)
}

func TakeScreenshot(ctx context.Context) ([]byte, error) {
	var buf []byte
	return buf, chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))
}

func TakeScreenshotOf(ctx context.Context, x, y, width, height float64) ([]byte, error) {
	log.Debug().Msgf("taking screenshot: [%f, %f] [%f, %f]", x, y, width, height)

	var buf []byte
	return buf, chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
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
	}))
}
