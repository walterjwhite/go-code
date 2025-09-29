package learning

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	completionUrl    = "https://www.linkedin.com/learning/me/my-library/completed"
	completedClasses = ".my-library"
)

func (s *Session) WasCompleted(course *Course) bool {
	log.Info().Str("course", course.Title).Msg("Session.WasCompleted - checking if course was completed")

	completionCtx, completionCancel := chromedp.NewContext(s.ctx)
	defer completionCancel()

	ctx, cancel := context.WithTimeout(completionCtx, *s.StepTimeout)
	defer cancel()

	err := action.Execute(ctx,
		chromedp.Navigate(completionUrl),
		chromedp.WaitReady(completedClasses),
		action.EndAction())
	if err != nil {
		log.Warn().Err(err).Msg("Session.WasCompleted - completion.0 - Error")
	}

	var nodes []*cdp.Node
	err = chromedp.Run(completionCtx,
		chromedp.Nodes(".completed-body", &nodes),
	)
	if err != nil {
		log.Warn().Err(err).Msg("Session.WasCompleted - completion.1 - Error")
		return false
	}

	if len(nodes) != 1 {
		log.Warn().Int("nodes", len(nodes)).Msg("Session.WasCompleted - expected only 1 matching element")
		return false
	}

	err = chromedp.Run(completionCtx,
		chromedp.Nodes(".completed-body .entity-link", &nodes),
	)
	if err != nil {
		log.Warn().Err(err).Msg("Session.WasCompleted - completion.2 - Error")
		return false
	}

	if len(nodes) == 0 {
		return false
	}

	log.Info().Int("nodes", len(nodes)).Msg("Session.WasCompleted - found completed courses")
	for _, node := range nodes {
		href, exists := node.Attribute("href")
		if !exists {
			log.Warn().Msg("Session.WasCompleted - href does not exist")
			continue
		}

		if course.Url == href {
			return true
		}
	}

	return false
}
