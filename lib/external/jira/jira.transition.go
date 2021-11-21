package jira

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"strconv"
)

func (i *Instance) Transition(ticketId string, transitionAction string) {
	i.setupAuth()

	_, err := i.client.Issue.DoTransition(ticketId, strconv.Itoa(i.TransitionActionMapping[transitionAction]))
	logging.Panic(err)
}
