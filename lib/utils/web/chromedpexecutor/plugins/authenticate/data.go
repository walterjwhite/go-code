package authenticate

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/time/delay"
	"github.com/walterjwhite/go-code/lib/time/keep_alive"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

type Session struct {
	ctx         context.Context
	Credentials *Credentials
	Website     *Website

	VisibleTimeout *time.Duration
	LocateDelay    delay.Delayer

	MinLocateDelay     *time.Duration
	DeviateLocateDelay *time.Duration

	chromedpsession session.ChromeDPSession
	activityChannel chan bool
	keepAlive       *keep_alive.KeepAlive
}

func (s *Session) With(ctx context.Context, chromedpsession session.ChromeDPSession) *Session {
	s.ctx = ctx
	s.chromedpsession = chromedpsession
	s.activityChannel = make(chan bool)
	if s.Website.SessionTimeout != nil {
		s.keepAlive = keep_alive.New(s.ctx, *s.Website.SessionTimeout, s.doKeepAlive)
	}

	if s.MinLocateDelay != nil {
		log.Debug().Msgf("delay: %v", *s.MinLocateDelay)
		if s.DeviateLocateDelay != nil {
			log.Debug().Msgf("deviate delay: %v", *s.DeviateLocateDelay)
			s.LocateDelay = delay.NewRandom(*s.MinLocateDelay, *s.DeviateLocateDelay)
		} else {
			s.LocateDelay = delay.New(*s.MinLocateDelay)
		}
	}

	return s
}

type Credentials struct {
	Secrets []*FieldSecret
}

type FieldSecret struct {
	FieldId     *string
	SecretKey   *string
	SecretValue *string
}

type Website struct {
	Url            *string
	FieldGroups    []*FieldGroup
	SessionTimeout *time.Duration

	SupportsMultipleTabs bool


	keepAliveActions []chromedp.Action
}

type FieldGroup struct {
	Fields         []*Field
	SubmitSelector *string
}

type Field struct {
	Id *string

	Selector *string
}

