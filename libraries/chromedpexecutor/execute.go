package chromedpexecutor

import (
	"context"
	"errors"
	"flag"
	"github.com/chromedp/chromedp"
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
	if len(*devToolsWsUrlFlag) == 0 {
		logging.Panic(errors.New("Please specify the DevToolsWSUrl"))
	}

	actxt, cancelActxt := chromedp.NewRemoteAllocator(ctx, *devToolsWsUrlFlag)

	// create new tab
	ctx, cancelCtxt := chromedp.NewContext(actxt)

	return &ChromeDPSession{Context: ctx, CancelAllocator: cancelActxt, CancelContext: cancelCtxt, Waiter: waiter}
}

func (s *ChromeDPSession) Execute(actions ...chromedp.Action) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)
		logging.Panic(chromedp.Run(s.Context, action))

		s.doWait(i, actions...)
	}
}

func (s *ChromeDPSession) doWait(i int, actions ...chromedp.Action) {
	if s.Waiter != nil {
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
