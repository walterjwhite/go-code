package learning

import (
	"fmt"
	"github.com/chromedp/chromedp"

	"context"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"

	"strings"
)

const (
	codeChallenge = ".classroom-code-challenge__frame"

	getSectionTitleQuery = `document.getElementsByClassName('classroom-nav__subtitle')[0].textContent.trim()`
	getSectionIndexQuery = `
		Array.from(document.getElementsByClassName('classroom-toc-item__title'))
		.findIndex(element => element.textContent.includes('%s'))
	`
	clickSectionQuery = `document.getElementsByClassName('classroom-toc-item__title')[%d].click()`
)

func (s *Session) isCoderPad() bool {
	return action.ExistsByCssSelector(s.ctx, codeChallenge)
}

func (s *Session) coderPadNext() error {
	title, err := s.getSectionTitle()
	if err != nil {
		logging.Warn(err, false, "coderPadNext.getSectionTitle")
		return err
	}

	index, err := s.getSectionIndex(title)
	if err != nil {
		logging.Warn(err, false, "coderPadNext.getSectionIndex")
		return err
	}

	return s.clickSection(index + 1)
}

func (s *Session) getSectionTitle() (string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var value string
	err := chromedp.Run(ctx,
		chromedp.Evaluate(getSectionTitleQuery, &value))

	log.Debug().Msgf("title: %s", value)
	return value, err
}

func (s *Session) getSectionIndex(title string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var value int

	return value, chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(getSectionIndexQuery, strings.ReplaceAll(title, "'", "\\'")), &value))
}

func (s *Session) clickSection(sectionIndex int) error {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(clickSectionQuery, sectionIndex), nil))
}
