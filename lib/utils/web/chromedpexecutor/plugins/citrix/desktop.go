package citrix

import (
  "context"
  "fmt"
  "github.com/chromedp/cdproto/target"

  "github.com/chromedp/chromedp"

  "github.com/avast/retry-go"
  "github.com/rs/zerolog/log"
  "github.com/walterjwhite/go-code/lib/application/logging"
  "github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
  "github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/citrix/plugins/mouse_wiggle"
  "github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"

  "sync"
  "time"
)


type CitrixWorker interface {
  Work(ctx context.Context, headless bool)
  Cleanup()
}

func (s *Session) Work() {
  waitGroup := &sync.WaitGroup{}

  for index := range s.Instances {
    log.Info().Msgf("Work [%v]", s.Instances[index])
    s.Instances[index].session = s

    waitGroup.Add(1)
    go s.Instances[index].run(waitGroup)
  }

  waitGroup.Wait()
  log.Info().Msg("Work() done")
}

func (s *Session) OnBreak(breakDuration *time.Duration) {
  log.Info().Msgf("taking break: %v", breakDuration)

  s.cleanup()
}

func (s *Session) OnStop() {
  s.Cancel()
}

func (i *Instance) run(waitGroup *sync.WaitGroup) {
  log.Info().Msgf("run.start [%v]", i)

  defer waitGroup.Done()

  log.Info().Msg("run")
  if !i.initialized {
    i.launch()

    movementWaitTime := 3 * time.Minute
    timeBetweenActions := 3 * time.Second
    i.Worker = &mouse_wiggle.State{MovementWaitTime: &movementWaitTime, TimeBetweenActions: &timeBetweenActions}

    i.initialized = true

    i.handlePromptStatic()

    if !i.actionsInitialized {
      log.Info().Msg("running actions")
      i.actions()
      i.actionsInitialized = true
    }
  }

  i.Worker.Work(i.ctx, i.session.Headless)
  log.Info().Msgf("run.end [%v]", i)
}

func (i *Instance) launch() {
  i.session.handleExpired()

  err := retry.Do(
    func() error {
      return i.tryLaunch()
    },
    retry.Attempts(3),
    retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
      return retry.BackOffDelay(n, err, config)
    }),
    retry.Delay(5*time.Second),
  )

  logging.Panic(err)
}

func (i *Instance) tryLaunch() error {
  log.Info().Msgf("launching instance: %d -> %d @ [%s]", i, i.Index, action.Location(i.session.ctx))
  targetElementXpath := fmt.Sprintf("//*[@id=\"home-screen\"]/div[2]/section[5]/div[5]/div/ul/li[%d]/a[1]", i.Index)
  targetIDChannel := chromedp.WaitNewTarget(i.session.ctx, matchTabWithNonEmptyURL)

  timeLimitedCtx, timeLimitedCancel := context.WithTimeout(i.session.ctx, 5*time.Second)
  defer timeLimitedCancel()

  log.Info().Msgf("clicking: %s", targetElementXpath)
  err := chromedp.Run(timeLimitedCtx, chromedp.Click(targetElementXpath))
  if err != nil {
    return err
  }

  log.Info().Msg("clicked")
  newInstance, newCancelFunc := chromedp.NewContext(i.session.ctx, chromedp.WithTargetID(<-targetIDChannel))
  i.ctx = newInstance
  i.cancel = newCancelFunc

  log.Info().Msg("new instance")
  return nil
}

func matchTabWithNonEmptyURL(info *target.Info) bool {
  return info.URL != ""
}

func (i *Instance) cleanup() {
  if !i.initialized {
    return
  }

  i.initialized = false
  i.Worker.Cleanup()

  i.cancel()
  i.cancel = nil
  i.ctx = nil

  log.Info().Msgf("cleanup - %d", i.Index)
}

func (i *Instance) actions() {
  if len(i.Actions) > 0 {
    log.Info().Msgf("running actions - delay: %v", *i.InitialActionDelay)
    time.Sleep(*i.InitialActionDelay)

    log.Info().Msgf("running actions: %v", i.Actions)
    logging.Panic(chromedp.Run(i.ctx, run.ParseActions(i.Actions...)...))
  }
}
