package monitor

var ConstructorRegistry = map[string]func(action *Action, session *Session) Monitor{}
var AlertRegistry = map[string]func(action *Action, session *Session, notification *notify.Notification) notify.Notifier{}
