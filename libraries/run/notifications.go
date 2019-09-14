package run

/*

import (
	"fmt"
	"log"
	"libraries/notify"
)

func buildErrorNotification(application string, err error, notificationBuilder func(notification notify.Notification) notify.Notifier) notify.Notifier {
	log.Printf("Error: %v\n", err)
	return notificationBuilder(notify.Notification{Id: "run", Title: fmt.Sprintf("%v Error", application), Details: fmt.Sprintf("%v", err)})
}

func checkIfStarted(application string, channel chan *string, notificationBuilder func(notification notify.Notification) notify.Notifier {
	for {
		select {
			case applicationStartedLine := <-channel:
			log.Printf("Application Started: %v\n", applicationStartedLine)
			notifyApplicationStarted(application, applicationStartedLine, notificationBuilder)
			return
		}
	}
}

func notifyApplicationStarted(application string, applicationStartedLine *string, notificationBuilder func(notification notify.Notification) notify.Notifier {
	return notificationBuilder(notify.Notification{Id: "run", Title: fmt.Sprintf("%v Started", application), Details: *applicationStartedLine})
}
*/
