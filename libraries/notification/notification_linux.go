package notification

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"os/exec"
)

type linuxNotification struct{}

func (n *linuxNotification) Notify(notification Notification) {
	args := []string{}

	if notification.Icon != "" {
		args = append(args, "-i", notification.Icon)
	}

	args = append(args, notification.Title)
	args = append(args, notification.Description)

	cmd := exec.Command("notify-send", args...)

	logging.Panic(cmd.Run())
}

func New() Notifier {
	return &linuxNotification{}
}
