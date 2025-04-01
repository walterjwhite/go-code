package gateway

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"sync"
)

func (s *Session) Launch(waitGroup *sync.WaitGroup) {
	for i, instance := range s.Instances {
		log.Info().Msgf("Launching instance: %d -> %d", i, instance.Index)
		targetIDChannel := chromedp.WaitNewTarget(s.session.Context(), matchTabWithNonEmptyURL)

		logging.Panic(chromedp.Run(s.session.Context(), chromedp.Click(fmt.Sprintf("//*[@id=\"home-screen\"]/div[2]/section[5]/div[5]/div/ul/li[%d]/a[1]", instance.Index))))

		newInstance, newCancelFunc := chromedp.NewContext(s.session.Context(), chromedp.WithTargetID(<-targetIDChannel))
		go s.RunInstance(newInstance, newCancelFunc, waitGroup, instance)
	}
}

func matchTabWithNonEmptyURL(info *target.Info) bool {
	return info.URL != ""
}

func (s *Session) RunInstance(instanceContext context.Context, instanceCancel context.CancelFunc, waitGroup *sync.WaitGroup, instance Instance) {
	waitGroup.Add(1)

	defer waitGroup.Done()
	defer instanceCancel()
	defer log.Info().Msgf("context done - %d", instance.Index)

	handlePrompt(instanceContext, instance)
	if log.Debug().Enabled() {
		chromedpexecutor.FullScreenshot(instanceContext, fmt.Sprintf("/tmp/2.gateway-prompt-%d.png", instance.Index))
	}

	s.wiggleMouse(instanceContext, instance)

	<-instanceContext.Done()
}
