package database

import (
	//"github.com/walterjwhite/go-application/libraries/database"
	"github.com/walterjwhite/go-application/libraries/io/yaml"
	"github.com/walterjwhite/go-application/libraries/utils/monitor"
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

	// TODO: this is configuration, this should be provided by another package
	yamlhelper.Read(action.Reference, &databaseMonitorAction)

	databaseMonitorAction.Action = action
	databaseMonitorAction.Session = session

	return databaseMonitorAction
}
