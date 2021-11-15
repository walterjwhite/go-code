package notification

import (
	"github.com/walterjwhite/go/lib/application"
	"github.com/walterjwhite/go/lib/application/logging"
)

type windowsNotification struct{}

func (n *windowsNotification) Notify(notification Notification) {
	toastNotification := toast.Notification{
		AppID:   application.GetApplicationId(),
		Title:   notification.Title,
		Message: notification.Description,

		/*
		   Actions: []toast.Action{
		       {"protocol", "I'm a button", ""},
		       {"protocol", "Me too!", ""},
		   },
		*/
	}

	if notification.Icon != "" {
		toastNotification.Icon = notification.Icon
	}

	logging.Panic(toastNotification.Push())
}

func New() Notifier {
	return &windowsNotification{}
}
