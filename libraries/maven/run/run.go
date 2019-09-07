package run

import (
	"context"
	"os/exec"
	//"libraries/notify"
)

func Run(ctx context.Context, profile string, debug bool /*, notificationBuilder func(notification notify.Notification) notify.Notifier*/) {
	var c Configuration
	c.getConf(profile)

	command := make([]exec.Cmd, 0)
	for index, application := range c.Applications {
		comands = append(commands, *runApplication(ctx, index, profile, c, application, debug, notificationBuilder))
	}

	for index, command := range commands {
		_, err := command.Process.Wait()
		if err != nil {
			//var notifier notify.Notifier
			//notifier = buildErrorNotification(c.Applications[index], err, notificationBuilder)
			//notifier.Notify()
		}
	}
}
