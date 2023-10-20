package notification

import (
	// "github.com/TheCreeper/go-notify"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/esiqveland/notify"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type bsdNotification struct{}

func (n *bsdNotification) Notify(notification Notification) {
	conn, err := dbus.SessionBusPrivate()
	logging.Panic(err)

	// ntf := notify.NewNotification(notification.Title, notification.Description)
	// ntf.AppIcon = notification.Icon
	// if len(notification.AudioFile) > 0 {
	// 	ntf.Hints = make(map[string]interface{})
	// 	ntf.Hints[notify.HintSoundFile] = notification.AudioFile
	// }

	// _, err := ntf.Show()

	dbusNotification := notify.Notification{
		// AppName:    "Test GO App",
		ReplacesID: uint32(0),
		// AppIcon:    iconName,
		Summary: notification.Title,
		Body:    notification.Description,
		Actions: []notify.Action{
			{Key: "cancel", Label: "Cancel"},
			{Key: "open", Label: "Open"},
		},
		Hints:         map[string]dbus.Variant{},
		ExpireTimeout: time.Second * 5,
	}

	_, err = notify.SendNotification(conn, dbusNotification)
	logging.Panic(err)
}

func New() Notifier {
	return &bsdNotification{}
}
