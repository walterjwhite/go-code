package monitor

import (
	"fmt"
)

type NotificationEvent struct {
	Session Session
	Action  Action
	//Notification notify.Notification
	Details string
}

func NewNotificationEvent(session *Session, action *Action, details string) *NotificationEvent {
	return *NotificationEvent{Session: *session, Action: *action, Details: details}
}

func getTitle(sessionDescription string, sessionActionDescription string, interval string) string {
	return fmt.Sprintf("%v / %v @ %v / %v\n", sessionDescription, sessionActionDescription, interval, timestamp.Get())
}

func (session *Session) watchChannel() {
	for {
		select {
		case notification, ok := <-session.Channel:
			if !ok {
				log.Printf("Channel appears to be shutdown\n")
				break
			}

			session.pushAlerts(notification)
		}
	}
}

func (session *Session) pushAlerts(notificationEvent *NotificationEvent) {
	session.scheduleNoActivityAlert()

	notification := notificationEvent.build()
	session.push(notification)
}

/*
func (n *NotificationEvent) build() notify.Notification {
}
*/

func (session *Session) push(notificationEvent *NotificationEvent) {
}

func (session *Session) NoActivityAlert() error {
	details := fmt.Sprintf(session.NoActivity.Description, session.NoActivity.Interval)
	// generate notification / push to channel

	go session.push(notification)
	session.scheduleNoActivityAlert()

	return nil
}
