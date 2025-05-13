package provider

import (
  "context"

  "github.com/chromedp/chromedp"
  "github.com/rs/zerolog/log"
  "github.com/walterjwhite/go-code/lib/application/logging"
)

type Conf struct {
  Headless bool
  ProxyAddress string
}

func New(c *Conf, ctx context.Context) (context.Context, context.CancelFunc) {
  ctx1, cancel := chromedp.NewContext(ctx)

  var opts []chromedp.ExecAllocatorOption

  if c.Headless {
    log.Info().Msgf("enabling headless")
    opts = append(opts, chromedp.Flag("headless", true))

    logging.Panic(chromedp.Run(ctx1, chromedp.EmulateViewport(1920, 1080)))
  }

  if len(c.ProxyAddress) > 0 {
    log.Info().Msgf("using proxy: %v", c.ProxyAddress)
    opts = append(opts, chromedp.ProxyServer(c.ProxyAddress))
  }

  ctx1, _ = chromedp.NewExecAllocator(ctx1, opts...)

  ctx1, _ = chromedp.NewContext(ctx1)
  return ctx1, cancel
}
