package learning

import (
	"github.com/chromedp/chromedp"
	"errors"
	"context"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

const (
	classroomQuiz = ".classroom-quiz"
	classroomNext = ".classroom-next-up__image-container"
)

func (s *Session) isQuiz() bool {
	return action.ExistsByCssSelector(s.ctx, classroomQuiz)
}

func (s *Session) quizNext() error {
	if action.ExistsByCssSelector(s.ctx, classroomNext) {
		ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
		defer cancel()

		return chromedp.Run(ctx, chromedp.Click(classroomNext))
	}

	return errors.New("classroomNext button not found")
}
