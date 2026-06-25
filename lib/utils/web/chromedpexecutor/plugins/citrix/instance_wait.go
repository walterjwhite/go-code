package citrix

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	citrixReady        = "GBufferFrameProducer(Main): ossSupportState is: false"
	citrixSessionReady = "Received session info :  {\"caption\":"


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
		chromedp.Poll(fmt.Sprintf(citrixPollStatement, strings.ReplaceAll(logMessage, "'", "\\'")), nil, chromedp.WithPollingInterval(pollInterval)))
}

func (i *Instance) waitForSessionReady() error {
	ctx, cancel := context.WithTimeout(i.ctx, citrixSessionInitializationTimeout)
	defer cancel()

	return waitForCitrixSessionReady(ctx)
}

func (i *Instance) waitForInitialization(pctx context.Context) error {
	ctx, cancel := context.WithTimeout(pctx, citrixSessionInitializationTimeout)
	defer cancel()

	return waitForCitrixInitialization(ctx)
}





func (i *Instance) waitForInitDelay() {
	if i.InitializationDelay > 0 {
		log.Info().Msgf("%v - Instance.init - delay - %v", i, i.InitializationDelay)
		time.Sleep(i.InitializationDelay)
	}
}
