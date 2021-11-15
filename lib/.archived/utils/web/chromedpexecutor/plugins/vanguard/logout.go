package vanguard

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

const (
	logoutButton = "//*[@id=\"globalNavUtilityBar\"]/div/div/ul/li[5]/a/span"
)

func (s *Session) Logout() {
	log.Info().Msg("Logging out")

	defer s.chromedpsession.Cancel()

	s.chromedpsession.Execute(chromedp.Click(logoutButton))
}
