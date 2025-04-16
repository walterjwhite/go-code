package citrix

import (
	"context"
	"github.com/walterjwhite/go-code/lib/time/periodic"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"sync"
	"time"
)

type Session struct {
	Credentials *Credentials
	Endpoint    *Endpoint

	Tickle   *Tickle
	Headless bool

	UseLightVersion bool

	Delay     *time.Duration
	Instances []Instance

	PostAuthenticationDelay   *time.Duration
	PostAuthenticationActions []string

	EndHour   int
	StartHour int

	LunchBreakStartHour int
	LunchBreakEndHour   int

	Timeout *time.Duration

	session session.ChromeDPSession

	keepAliveChannel <-chan time.Time
	waitGroup        sync.WaitGroup
}

func (s *Session) Session() session.ChromeDPSession {
	return s.session
}

type Credentials struct {
	Domain   string
	Username string
	Password string

	Pin string
}

type Endpoint struct {
	Uri string

	UsernameXPath    string
	PasswordXPath    string
	TokenXPath       string
	LoginButtonXPath string

	AuthenticationDelay *time.Duration
}

type Tickle struct {
	TickleInterval   *time.Duration
	periodicInstance *periodic.PeriodicInstance
}

type Instance struct {
	Index  int
	Action string

	InitialActionDelay *time.Duration
	TimeBetweenActions *time.Duration

	Actions []string

	MovementWaitTime *time.Duration

	Pomodoro *PomodoroInstance

	ctx    context.Context
	cancel context.CancelFunc

	session *Session

	lastMouseX float64
	lastMouseY float64

	initialized        bool
	actionsInitialized bool

	breakChannel chan *time.Duration
	stopChannel  chan bool
}
