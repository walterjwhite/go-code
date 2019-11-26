package chromedpexecutor

import (
	"context"
	"flag"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/sleep"
	"time"
)

type ChromeDPSession struct {
	Context context.Context
	Waiter  sleep.Waiter

	CancelAllocator context.CancelFunc
	CancelContext   context.CancelFunc
}

type TimeLimitedChromeAction struct {
	Action      chromedp.Action
	Limit       time.Duration
	IsException bool
}

var (
	devToolsWsUrlFlag = flag.String("DevToolsWSUrl", "", "Dev Tools WS URL")

	// TODO: add flags to tweak the deviation and minimum wait times
	// OR if a fixed delay is preferred
	waiter sleep.Waiter
)

func init() {
	waiter = &sleep.RandomDelay{MinimumDelay: 250, Deviation: 5000}
}

func New(ctx context.Context) *ChromeDPSession {
	actxt, cancelActxt := chromedp.NewRemoteAllocator(ctx, *devToolsWsUrlFlag)

	// create new tab
	ctx, cancelCtxt := chromedp.NewContext(actxt)

	return &ChromeDPSession{Context: ctx, CancelAllocator: cancelActxt, CancelContext: cancelCtxt, Waiter: waiter}
}

func (s *ChromeDPSession) Execute(actions ...chromedp.Action) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)
		logging.Panic(chromedp.Run(s.Context, action))

		if i < (len(actions) - 1) {
			s.Waiter.Wait()
		}
	}
}

func (s *ChromeDPSession) ExecuteTimeLimited(actions ...TimeLimitedChromeAction) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)

		ctx, cancel := context.WithTimeout(s.Context, action.Limit)
		defer cancel()

		logging.Warn(chromedp.Run(ctx, action.Action), action.IsException)

		if i < (len(actions) - 1) {
			s.Waiter.Wait()
		}
	}
}

func (s *ChromeDPSession) Cancel() {
	s.CancelAllocator()
	s.CancelContext()
}

// TODO: this does not work and just hangs
func (s *ChromeDPSession) Screenshot(filename string) {
	var buf []byte

	log.Info().Msgf("capturing screenshot: %v", filename)
	//logging.Panic(chromedp.Run(s.Context, fullScreenshot(90, &buf)))
	s.Execute(chromedp.CaptureScreenshot(&buf))

	log.Info().Msgf("took screenshot - writing to: %v", filename)

	logging.Panic(ioutil.WriteFile(filename, buf, 0644))
	log.Info().Msgf("captured screenshot: %v", filename)
}

func fullScreenshot(quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
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
