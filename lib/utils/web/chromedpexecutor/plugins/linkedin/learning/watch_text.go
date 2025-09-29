package learning

import (
	"github.com/chromedp/chromedp"

	"context"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	classroomTextPaging = ".classroom-multimedia__paging"

	classroomTextNextClick = `(function() {
		Array.from(document.querySelectorAll('button'))
			.find(d => d.textContent.trim() === 'Next').click();

		e = Array.from(document.querySelectorAll('button'))
			.find(d => d.textContent.trim() === 'Next');
		if(e !== undefined) {
			e.click();
			return true;
		}

		return false;
	})()
	`
)

func (s *Session) isText() bool {
	return action.ExistsByCssSelector(s.ctx, classroomTextPaging)
}

func (s *Session) textNext() (bool, error) {
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	var exists bool
	return exists, chromedp.Run(ctx,
		chromedp.Evaluate(classroomTextNextClick, &exists))
}
