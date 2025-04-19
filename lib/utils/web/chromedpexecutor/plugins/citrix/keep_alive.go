package citrix

import (
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"strings"
)

func (s *Session) keepAlive() {
	for {
		select {
		case <-s.keepAliveChannel:
			s.doKeepAlive()
		case <-s.ctx.Done():
			log.Warn().Msg("session context ended, exiting keep-alive")
			return
		case <-application.Context.Done():
			log.Warn().Msg("application context ended, exiting keep-alive")
			return
		}
	}
}

func (s *Session) doKeepAlive() {
	s.handleExpired()

	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	log.Debug().Msgf("tickling: %v", s.Endpoint.Uri)
	action.Execute(s.ctx, chromedp.Navigate(s.Endpoint.Uri))
}

func (s *Session) handleExpired() {
	if s.isExpired() {
		logging.Panic(errors.New("session expired"))
	}
}

func (s *Session) isExpired() bool {
	currentUrl := action.Location(s.ctx)
	if strings.HasSuffix(currentUrl, "/logout.html") {
		return true
	}

	return strings.HasSuffix(currentUrl, "LogonPoint/tmindex.html")
}
