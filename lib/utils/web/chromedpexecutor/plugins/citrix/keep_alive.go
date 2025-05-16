package citrix

import (
  "context"
  "errors"
  "github.com/avast/retry-go"
  "github.com/chromedp/chromedp"
  "github.com/rs/zerolog/log"
  "github.com/walterjwhite/go-code/lib/application"
  "github.com/walterjwhite/go-code/lib/application/logging"
  "github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
  "strings"
  "time"
)

func (s *Session) keepAlive() {
  for {
    select {
    case <-s.keepAliveChannel:
      s.doKeepAlive()
    case <-s.ctx.Done():
      log.Warn().Msg("session context ended, exiting keep-alive")
      return
    case <-application.Context.Done():
      log.Warn().Msg("application context ended, exiting keep-alive")
      return
    }
  }
}

func (s *Session) doKeepAlive() {
  start := time.Now()

  s.handleExpired()

  err := retry.Do(
    func() error {
      return s.doTryKeepAlive()
    },
    retry.Attempts(5),
    retry.Delay(3*time.Second),
  )

  end := time.Now()
  delta := end.Sub(start)
  log.Warn().Msgf("Delta: %v", delta)

  if err != nil {
    action.FullScreenshot(s.ctx, "/tmp/citrix-keep-alive-timeout-error.png")
  } else {
    action.FullScreenshot(s.ctx, "/tmp/citrix-keep-alive.png")
  }

  logging.Panic(err)
}

func (s *Session) doTryKeepAlive() error {
  ctx, cancel := context.WithTimeout(s.ctx, *s.KeepAliveTimeout)
  defer cancel()

  log.Debug().Msgf("tickling: %v", s.Endpoint.Uri)
  return chromedp.Run(ctx, chromedp.Navigate(s.Endpoint.Uri))
}

func (s *Session) handleExpired() {
  if s.isExpired() {
    logging.Panic(errors.New("session expired"))
  }
}

func (s *Session) isExpired() bool {
  currentUrl := action.Location(s.ctx)
  if strings.HasSuffix(currentUrl, "/logout.html") {
    return true
  }

  return strings.HasSuffix(currentUrl, "LogonPoint/tmindex.html")
}
