package citrix

import (
	"errors"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"strings"
)

func (s *Session) keepAlive() {
	defer s.waitGroup.Done()

	for range s.keepAliveChannel {
		sessionMutex.Lock()
		defer sessionMutex.Unlock()

		var currentUrl string
		logging.Panic(chromedp.Run(s.session.Context(), chromedp.Location(&currentUrl)))
		if strings.HasSuffix(currentUrl, "/logout.html") {
			logging.Panic(errors.New("server ended the session - logout.html"))
		} else if strings.HasSuffix(currentUrl, "LogonPoint/tmindex.html") {
			logging.Panic(errors.New("server ended the session - tmindex.html"))
		}

		log.Debug().Msgf("tickling: %v", s.Endpoint.Uri)
		session.Execute(s.session, chromedp.Navigate(s.Endpoint.Uri))
	}
}
