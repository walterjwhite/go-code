package learning

import (
	"context"

	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

const (
	skipSurveyText = "Skip survey"
)

func (s *Session) isSurvey() bool {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(clickElementWithText, skipSurveyText), &exists))

	logging.Warn(err, false, "isSurvey")
	return exists
}
