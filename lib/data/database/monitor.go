package database

import (
	//"github.com/walterjwhite/go-code/lib/database"
	"github.com/walterjwhite/go-code/lib/io/yaml"
	"github.com/walterjwhite/go-code/lib/utils/monitor"
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

	yaml.Read(action.Reference, &databaseMonitorAction)

	databaseMonitorAction.Action = action
	databaseMonitorAction.Session = session

	return databaseMonitorAction
}
