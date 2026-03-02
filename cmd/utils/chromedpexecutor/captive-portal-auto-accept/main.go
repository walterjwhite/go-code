package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/provider"

	"github.com/chromedp/chromedp"

	"path/filepath"
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

	if err := validateURL(captivePortalSession.Url); err != nil {
		logging.Error(err)
		return
	}

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

	filename := fmt.Sprintf("%d.connectivity-check.png", index)
	safePath := filepath.Join("/tmp", filepath.Base(filename))
	action.Screenshot(captivePortalSession.ctx, safePath)
}

func validateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL is required and cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme '%s': only http and https are allowed", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a valid hostname")
	}

	return nil
}
