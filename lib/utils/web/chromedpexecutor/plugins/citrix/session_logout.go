package citrix

import (
	"context"

	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	logoutTimeout = 5 * time.Second

	menuButtonXpath   = "//*[@id=\"userMenuBtn\"]/div"
	logoffButtonXpath = "//*[@id=\"menuLogOffBtn\"]"
)

func (s *Session) Logout() error {
	ctx, cancel := context.WithTimeout(s.ctx, logoutTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.Click(menuButtonXpath),
		chromedp.Click(logoffButtonXpath),
	)
}
