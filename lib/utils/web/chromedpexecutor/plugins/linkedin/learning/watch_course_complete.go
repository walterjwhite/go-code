package learning

import (
	"context"

	"fmt"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"strings"
)

func (s *Session) isCourseComplete(course *Course) bool {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var courseComplete bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(hasElementWithText, strings.ReplaceAll(course.Title, "'", "\\'")), &courseComplete))

	logging.Warn(err, false, "isCourseComplete")
	return courseComplete
}
