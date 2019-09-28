package monitor

import (
	"context"
	"log"
)

type Monitor interface {
	Execute()
}

func New(ctx context.Context, configurationFile string) *Session {
	session := &Session{}
	session.Channel = make(chan *NotificationEvent)

	session.setup(ctx, configurationFile)
	session.setupActions()
	session.scheduleNoActivityAlert()

	return session
}

func (session *Session) setup(ctx context.Context, configurationFile string) {
	log.Printf("session: %v\n", session)
	read(configurationFile, session)
	session.Context = ctx
}

func (session *Session) setupActions() {
	log.Printf("actions: %v\n", session.Actions)

	for _, action := range session.Actions {
		session.setupAction(action)
	}

	go session.watchChannel()
}

func (session *Session) setupAction(action Action) {
	action.Session = session

	log.Printf("Action: %v\n", action.Type)
	log.Printf("ConstructorRegistry: %v\n", ConstructorRegistry)

	monitorAction := ConstructorRegistry[action.Type](&action, session)
	action.Monitor = monitorAction

	log.Printf("Monitor: %v\n", monitorAction)
	action.schedule()
}
