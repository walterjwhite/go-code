package provider

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

func (c *Conf) getAllocator(pctx context.Context) (context.Context, context.CancelFunc) {
	if len(c.Remote) == 0 {
		log.Info().Msgf("using exec: %v", c.Remote)
		return chromedp.NewExecAllocator(pctx, c.getOptions()...)
	}

	log.Info().Msgf("using remote: %v", c.Remote)
	return chromedp.NewRemoteAllocator(pctx, c.Remote)
}

func (c *Conf) getOptions() []chromedp.ExecAllocatorOption {
	var opts []chromedp.ExecAllocatorOption



	if c.Headless {
		if len(c.Remote) > 0 {
			log.Warn().Msg("conflicting options, cannot enable headless with remote")
		} else {
			log.Info().Msgf("enabling headless")
			opts = append(opts, chromedp.Flag("headless", true))
		}
	}

	if len(c.ProxyAddress) > 0 {
		log.Info().Msgf("using proxy: %v", c.ProxyAddress)
		opts = append(opts, chromedp.ProxyServer(c.ProxyAddress))
	}

	return opts
}
