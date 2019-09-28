package monitor

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/timestamp"
	"log"
)

type NotificationEvent struct {
	Session Session
	Action  Action
	//Notification notify.Notification
	Details string
}

func NewNotificationEvent(session *Session, action *Action, details string) *NotificationEvent {
	return &NotificationEvent{Session: *session, Action: *action, Details: details}
}

func GetTitle(sessionDescription string, sessionActionDescription string, interval string) string {
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

	//notification := notificationEvent.build()
	// notificationEvent should be notification (as it will be sent to Windows10 (toast), dbus, etc.
	session.push(notificationEvent)
}

/*
func (n *NotificationEvent) build() notify.Notification {
}
*/

func (session *Session) push(notificationEvent *NotificationEvent) {
	/*
		for _, alert := range session.Alerts {
				notifier := AlertRegistry[alert.Type](&alert, session, &notification)
				notifier.Notify()
		}
	*/
}

func (session *Session) NoActivityAlert() error {
	details := fmt.Sprintf(session.NoActivity.Description, session.NoActivity.Interval)
	// generate notification / push to channel

	notificationEvent := &NotificationEvent{Session: *session, Details: details}

	go session.push(notificationEvent)
	session.scheduleNoActivityAlert()

	return nil
}
