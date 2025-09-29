package learning

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	inProgressXPath = `//*[@id="hue-tabs-ember58-tab-me.my-library.in-progress"]/span`
)

func (s *Session) InProgress() []*Course {
	log.Info().Msg("Session.InProgress - start")

	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	err := action.Execute(ctx,
		chromedp.Navigate(linkedInLearningUrl+"/me/my-library/in-progress"),
		chromedp.WaitReady(inProgressXPath))
	if err != nil {
		log.Warn().Err(err).Msg("Session.InProgress - error waiting for in progress to load")
		return nil
	}

	return s.extractCourses()
}
