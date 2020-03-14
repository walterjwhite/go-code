package jira

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"strconv"
)

func (i *Instance) Transition(ticketId string, transitionAction string) {
	_, err := i.client.Issue.DoTransition(ticketId, strconv.Itoa(i.TransitionActionMapping[transitionAction]))
	logging.Panic(err)
}
