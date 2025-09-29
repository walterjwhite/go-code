package learning

import (
	"context"

	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

const (
	courseCompleteXPath = "media-screen__content media-screens-course-complete__content"
)

func (s *Session) isCourseComplete() bool {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var courseComplete bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(hasElementWithText, courseCompleteXPath), &courseComplete))

	logging.Warn(err, false, "isCourseComplete")
	return courseComplete
}
