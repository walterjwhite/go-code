package database

import (
	//"github.com/walterjwhite/go-application/libraries/database"
	"github.com/walterjwhite/go-application/libraries/monitor"
	"github.com/walterjwhite/go-application/libraries/yamlhelper"
)

type DatabaseMonitorAction struct {
	Description string

	Query Query

	Action  *monitor.Action
	Session *monitor.Session
}

type DatabaseMonitorActionEvent struct {
	DatabaseMonitorAction *DatabaseMonitorAction
	Event                 interface{}
}

func NewMonitor(action *monitor.Action, session *monitor.Session) DatabaseMonitorAction {
	var databaseMonitorAction DatabaseMonitorAction

	yamlhelper.Read(action.Reference, &databaseMonitorAction)

	databaseMonitorAction.Action = action
	databaseMonitorAction.Session = session

	return databaseMonitorAction
}
