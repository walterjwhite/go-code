package notification

import (
	"github.com/walterjwhite/go-application/libraries/identifier"
	"github.com/walterjwhite/go-application/libraries/logging"
	"gopkg.in/toast.v1"
)

type windowsNotification struct{}

func (n *windowsNotification) Notify(notification Notification) {
	toastNotification := toast.Notification{
		AppID:   identifier.GetApplicationId(),
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
