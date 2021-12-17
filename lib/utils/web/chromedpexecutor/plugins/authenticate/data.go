package authenticate

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/time/keep_alive"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

type Session struct {
	ctx         context.Context
	Credentials *Credentials
	Website     *Website

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

	return s
}

type Credentials struct {
	Secrets []*FieldSecret
}

type FieldSecret struct {
	// must match Field.id to be injected
	FieldId     *string
	SecretKey   *string
	SecretValue *string
}

type Website struct {
	Url            *string
	FieldGroups    []*FieldGroup
	SessionTimeout *time.Duration

	SupportsMultipleTabs bool

	//KeepAliveUrl   *string
	//IsKeepAlive bool

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

// TODO: add support for entering one-time tokens
// TODO: add support for answering challenge questions
