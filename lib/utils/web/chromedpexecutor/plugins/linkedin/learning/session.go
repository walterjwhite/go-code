package learning

import (
	"github.com/chromedp/cdproto/browser"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/publisher"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
)

func (s *Session) Run(publisher publisher.Publisher) {
	defer s.cancel()

	if len(s.EmailAddress) == 0 {
		log.Panic().Msg("Session.Run - email address is undefined")
	}
	if len(s.Password) == 0 {
		log.Panic().Msg("Session.Run - password is undefined")
	}

	action.Grant(s.ctx, []browser.PermissionType{})

	s.authenticate(publisher)

	for {
		s.consumeContent()
	}
}

func (s *Session) consumeContent() {
	log.Info().Msg("Session.consumeContent - start")

	courses := s.InProgress()
	if len(courses) == 0 {
		courses = s.Search(s.SearchCriteria)
	}

	for i := range courses {
		log.Info().Str("course", courses[i].Title).Str("url", courses[i].Url).Msg("Session.consumeContent - Found course")
		s.watch(courses[i])
	}

}
