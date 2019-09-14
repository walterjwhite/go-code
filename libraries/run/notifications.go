package run

import (
	"fmt"
	"github.com/walterjwhite/go-application/libraries/notification"
	"log"
)

// TODO: the channel is only sending the matching line
// perhaps we should instead send a notification on the channel
// then we can configure how we want to receive the notifications here, OS notification, email, sms, etc.
func monitorChannel(application string, channel chan *string) {
	for {
		select {
		case applicationStartedLine := <-channel:
			log.Printf("Application Started: %v\n", applicationStartedLine)
			notification.NotifierInstance.Notify(notification.Notification{Title: fmt.Sprintf("run: %v", application), Description: *applicationStartedLine, Type: notification.Info})
		}
	}
}
