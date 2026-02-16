package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"

	"github.com/chromedp/chromedp"

	"time"
)

type CaptivePortalSession struct {
	Url     string
	Actions []string

	Conf *provider.Conf

	ActionTimeout time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
}

var captivePortalSession = &CaptivePortalSession{}

func init() {
	application.Configure(captivePortalSession)

	if captivePortalSession.ActionTimeout <= 0 {
		captivePortalSession.ActionTimeout = time.Duration(5 * time.Second)
	}
}

func main() {
	defer application.OnPanic()
	captivePortalSession.ctx, captivePortalSession.cancel = provider.New(captivePortalSession.Conf, application.Context)
	defer captivePortalSession.cancel()

	runStep(0, chromedp.Navigate(captivePortalSession.Url))

	for i, action := range run.ParseActions(captivePortalSession.Actions...) {
		runStep(i+1, action)
	}

	log.Info().Msgf("sleeping %d(ns)", captivePortalSession.ActionTimeout)
	time.Sleep(captivePortalSession.ActionTimeout)

	runStep(len(captivePortalSession.Actions)+1, chromedp.Navigate(captivePortalSession.Url))
}

func runStep(index int, chromeAction chromedp.Action) {
	stepTimeoutContext, stepFetchCancel := context.WithTimeout(captivePortalSession.ctx, captivePortalSession.ActionTimeout)
	defer stepFetchCancel()
	logging.Error(chromedp.Run(stepTimeoutContext, chromeAction))

	action.Screenshot(captivePortalSession.ctx, fmt.Sprintf("/tmp/%d.connectivity-check.png", index))
}
