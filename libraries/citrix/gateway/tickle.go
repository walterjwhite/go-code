package gateway

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/periodic"
)

const (
	allAppsButtonXpath = "//*[@id=\"allAppsBtn\"]"
	desktopButtonXpath = "/html/body/section[3]/div[2]/header/section/a[2]/span"
)

// periodically click on the apps / desktops
func (s *Session) tickle(ctx context.Context) {
	if s.Tickle.periodicInstance != nil {
		s.Tickle.periodicInstance.Cancel()
		s.Tickle.periodicInstance = nil
	}

	// TODO: 1 set a timer and periodically call doTickle
	s.Tickle.periodicInstance = periodic.Periodic(ctx, s.Tickle.TickleInterval, s.doTickle)
}

func (s *Session) doTickle() error {
	// click a link periodically
	s.chromedpsession.Execute([]chromedp.Action{
		chromedp.Click(s.getAndSetTickleButton())}...)

	return nil
}

func (s *Session) getAndSetTickleButton() string {
	l := s.Tickle.lastTickledAll
	s.Tickle.lastTickledAll = !s.Tickle.lastTickledAll

	if l {
		return desktopButtonXpath
	}

	return allAppsButtonXpath
}
