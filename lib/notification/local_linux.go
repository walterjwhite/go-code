package notification

import (
	"github.com/TheCreeper/go-notify"
	"github.com/walterjwhite/go/lib/application/logging"
)

type linuxNotification struct{}

func (n *linuxNotification) Notify(notification Notification) {
	ntf := notify.NewNotification(notification.Title, notification.Description)
	ntf.AppIcon = notification.Icon
	if len(notification.AudioFile) > 0 {
		ntf.Hints = make(map[string]interface{})
		ntf.Hints[notify.HintSoundFile] = notification.AudioFile
	}

	_, err := ntf.Show()
	logging.Panic(err)
}

func New() Notifier {
	return &linuxNotification{}
}
