package run

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/notification"
	"github.com/walterjwhite/go-application/libraries/time/wait"
	"net"
)

func (a *Application) checkPort(ctx context.Context) {
	if a.Port <= 0 {
		return
	}

	p := &portCheck{port: a.Port}

	wait.Wait(ctx, &a.PortMonitorInterval, &a.PortMonitorTimeout, p.isPortOpen)

	m := fmt.Sprintf("Application Port Opened: %v\n", a.Port)
	log.Info().Msg(m)
	notification.NotifierInstance.Notify(notification.Notification{Title: fmt.Sprintf("run: %v", a.Name), Description: m, Type: notification.Info})
}

type portCheck struct {
	port int
}

func (c *portCheck) isPortOpen() bool {
	_, err := net.Dial("tcp", fmt.Sprintf(":%v", c.port))
	return err == nil
}
