package monitor

import (
	"github.com/walterjwhite/go-application/libraries/database"
)

type DatabaseMonitorAction struct {
	Description string

	Query database.Query

	Action  *monitor.Action
	Session *monitor.Session
}

type DatabaseMonitorActionEvent struct {
	DatabaseMonitorAction *DatabaseMonitorAction
	Event                 interface{}
}

func New(action *monitor.Action, session *monitor.Session) DatabaseMonitorAction {
	var databaseMonitorAction DatabaseMonitorAction

	yamlhelper.Read(action.Reference, &databaseMonitorAction)

	databaseMonitorAction.Action = action
	databaseMonitorAction.Session = session

	return databaseMonitorAction
}
