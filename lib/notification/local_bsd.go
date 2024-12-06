package notification

import (
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/esiqveland/notify"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type bsdNotification struct{}

func (n *bsdNotification) Notify(notification Notification) {
	conn, err := dbus.SessionBusPrivate()
	logging.Panic(err)



	dbusNotification := notify.Notification{
		ReplacesID: uint32(0),
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
