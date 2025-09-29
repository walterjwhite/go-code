package citrix

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	citrixReady = "GBufferFrameProducer(Main): ossSupportState is: false"

	citrixSessionReady = "Received session info :  {\"caption\":"

	citrixFatal = "Session launch failed - Citrix Workspace app cannot connect to the server"

	citrixPollStatement = `
	typeof globalThis.html5Logger !== 'undefined' && 
		globalThis.html5Logger.dump().some(log => typeof log === 'string' && log.includes('%s'))
	`

	citrixSessionInitializationTimeout = 30 * time.Second
	citrixPollInterval                 = 500 * time.Millisecond
)

func waitForCitrixInitialization(ctx context.Context) error {
	return waitForCitrixLog(ctx, citrixPollInterval, citrixReady)
}

func waitForCitrixSessionReady(ctx context.Context) error {
	return waitForCitrixLog(ctx, citrixPollInterval, citrixSessionReady)
}

func waitForCitrixLog(ctx context.Context, pollInterval time.Duration, logMessage string) error {
	log.Debug().Msgf("waitForCitrixLog - waiting for message: %s", logMessage)
	return chromedp.Run(ctx,
		chromedp.Poll(fmt.Sprintf(citrixPollStatement, logMessage), nil, chromedp.WithPollingInterval(pollInterval)))
}
