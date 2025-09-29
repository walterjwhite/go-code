package learning

import (
	"context"

	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

const (
	continueCourseText = "Continue course"
)

func (s *Session) isContinueCourse() bool {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(clickElementWithText, continueCourseText), &exists))

	logging.Warn(err, false, "Session.isContinueCourse")
	return exists
}
