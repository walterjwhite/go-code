package notification

import (
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/runner"

	"context"
	"fmt"
	"os/user"
	"strings"
)

type remoteNotification struct {
	Username string
	Server   string

	// acknowledge message, wait for seconds, 0 (do not wait)
	TimeWait uint

	Context context.Context
}

func (n *remoteNotification) Notify(notification Notification) {
	cmd := "msg"
	arguments := make([]string, 0)

	arguments = append(arguments, n.Username)
	arguments = append(arguments, fmt.Sprintf("/Server:%v", n.Server))

	if n.TimeWait > 0 {
		arguments = append(arguments, fmt.Sprintf("/Time:%v", n.TimeWait))
	}

	arguments = append(arguments, format(notification))

	_, err := runner.Run(n.Context, cmd, arguments...)
	logging.Panic(err)
}

func format(notification Notification) string {
	return /*notification.ApplicationId + ":" + */ notification.Title + ":" + notification.Description
}

func NewRemoteNotification(ctx context.Context, targetServer string, timeWait uint) Notifier {
	return &remoteNotification{Username: getUsername(), Server: targetServer, Context: ctx, TimeWait: timeWait}
}

func getUsername() string {
	currentUser, err := user.Current()
	logging.Panic(err)

	// remove domain
	return strings.Split(currentUser.Username, "\\")[1]
}
