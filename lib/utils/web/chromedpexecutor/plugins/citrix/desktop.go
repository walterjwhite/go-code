package citrix

import (
	"fmt"
	"github.com/chromedp/cdproto/target"

	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/plugins/run"
	"sync"
	"time"
)

var sessionMutex sync.Mutex

func (s *Session) Launch() {
	go s.scheduleEnd()
	s.waitUntilStart()

	for _, instance := range s.Instances {
		instance.session = s

		s.waitGroup.Add(1)
		go instance.run()
	}

	s.waitGroup.Wait()
}

func (i *Instance) run() {
	defer i.cleanup()

	i.breakChannel = make(chan *time.Duration)
	i.stopChannel = make(chan bool)

	defer close(i.breakChannel)
	defer close(i.stopChannel)

	go i.Pomodoro.init(i.breakChannel)
	go i.session.initLunchBreak(i.breakChannel)

	for {
		select {
		case <-i.stopChannel:
			log.Warn().Msgf("exiting instance: %v", i)
			return
		case duration := <-i.breakChannel:
			log.Warn().Msgf("taking a break: %v (%v)", i, *duration)

			i.cancel()
			i.initialized = false
			time.Sleep(*duration)
		default:
			log.Info().Msg("default case")

			if !i.initialized {
				log.Info().Msg("launching instance")
				i.launch()

				time.Sleep(*i.InitialActionDelay)
				i.handlePrompt()

				if !i.actionsInitialized {
					log.Info().Msg("running actions")
					i.actions()
					i.actionsInitialized = true
				}
			}

			i.wiggleMouse()

			time.Sleep(*i.TimeBetweenActions)
		}
	}

	log.Info().Msgf("end of for loop: %v", i)
}

func (i *Instance) waitForDesktop() {
	screenCheckDuration := 5 * time.Second
	if !i.isScreenLocked() {
		log.Warn().Msg("screen is not locked (windows icon is present)")
		return
	}

	for {
		if i.isWaitingForTermsAcceptance() {
			i.handlePrompt()
		} else if i.isLoggingIn() {
			log.Info().Msg("waiting for system to login")
			for {
				time.Sleep(screenCheckDuration)
				if !i.isLoggingIn() {
					log.Info().Msg("logged in")
					return
				}
			}
		} else {
			log.Warn().Msg("neither waiting for terms acceptance or logging in")
		}

		time.Sleep(screenCheckDuration)
	}
}

func (s *Session) waitUntilStart() {
	_ = waitUntil(s.StartHour)
}

func waitUntil(hour int) bool {
	currentTime := time.Now()
	if currentTime.Hour() < hour {
		targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hour, 0, 0, 0, currentTime.Location())
		duration := targetTime.Sub(currentTime)

		time.Sleep(duration)
		return true
	}

	return false
}

func (i *Instance) launch() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	log.Info().Msgf("Launching instance: %d -> %d", i, i.Index)
	targetIDChannel := chromedp.WaitNewTarget(i.session.session.Context(), matchTabWithNonEmptyURL)

	logging.Panic(chromedp.Run(i.session.session.Context(), chromedp.Click(fmt.Sprintf("//*[@id=\"home-screen\"]/div[2]/section[5]/div[5]/div/ul/li[%d]/a[1]", i.Index))))

	newInstance, newCancelFunc := chromedp.NewContext(i.session.session.Context(), chromedp.WithTargetID(<-targetIDChannel))
	i.ctx = newInstance
	i.cancel = newCancelFunc

	i.initialized = true
}

func matchTabWithNonEmptyURL(info *target.Info) bool {
	return info.URL != ""
}

func (i *Instance) cleanup() {
	i.session.waitGroup.Done()
	i.cancel()
	log.Info().Msgf("context done - %d", i.Index)
}

func (i *Instance) actions() {
	if len(i.Actions) > 0 {
		log.Info().Msgf("running actions - delay: %v", *i.InitialActionDelay)
		time.Sleep(*i.InitialActionDelay)

		log.Info().Msgf("running actions: %v", i.Actions)
		logging.Panic(chromedp.Run(i.ctx, run.ParseActions(i.Actions...)...))
	}
}
